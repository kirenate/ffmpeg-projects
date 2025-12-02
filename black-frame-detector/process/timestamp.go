package process

import (
	"encoding/json"
	"github.com/pkg/errors"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type VideoInfo struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
	Streams []struct {
		FrameRate string `json:"r_frame_rate"`
	} `json:"streams"`
}

func DurationToTimestamp(filename string, in <-chan string, wg *sync.WaitGroup) chan string {
	out := make(chan string)

	go func() {
		defer wg.Done()

		vinfo, err := GetVideoInfo(filename)
		if err != nil {
			log.Println("get video info", err.Error())
			return
		}

		framerate, err := strconv.Atoi(strings.Split(vinfo.Streams[0].FrameRate, "/")[0])
		if err != nil {
			log.Println("framerate not recognized", err.Error())
			return
		}

		for data := range in {
			d, err := strconv.Atoi(data)
			if err != nil {
				log.Println("strconv", err.Error())
				return
			}

			t, err := time.Parse("5", strconv.Itoa(d/framerate))
			if err != nil {
				log.Println("parse time", err.Error())
			}

			tm := t.Format("15:04:05")
			
			frames := strconv.Itoa(d % framerate)
			if len(frames) < 2 {
				frames = "0" + frames
			}

			tm += ":" + frames

			out <- tm
		}

		close(out)
	}()

	return out
}

func GetVideoInfo(filename string) (*VideoInfo, error) {
	kwargs := make(ffmpeg.KwArgs)
	kwargs["of"] = "json"
	kwargs["hide_banner"] = ""
	data, err := ffmpeg.Probe(filename, kwargs)
	if err != nil {
		return nil, errors.Wrap(err, "Error probing, file might be corrupted")
	}

	var videoInfo VideoInfo
	err = json.Unmarshal([]byte(data), &videoInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Duration not found! Corrupted file!")
	}

	return &videoInfo, nil
}
