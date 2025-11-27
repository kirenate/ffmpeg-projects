package process

import (
	"encoding/json"
	"github.com/pkg/errors"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"strconv"
	"strings"
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

func DurationToTimestamp(vinfo *VideoInfo) (string, error) {
	sl := strings.Split(vinfo.Format.Duration, ".")

	rem, err := strconv.Atoi(sl[1][0:2])
	if err != nil {
		return "", errors.Wrap(err, "failed to convert rem from string to int")
	}

	framerate, err := strconv.Atoi(strings.Split(vinfo.Streams[0].FrameRate, "/")[0])
	if err != nil {
		return "", errors.Wrap(err, "framerate not recognized")
	}

	t, err := time.Parse("5", sl[0])
	if err != nil {
		return "", errors.Wrap(err, "parse time")
	}
	for rem >= framerate {
		t.Add(1 * time.Second)
		rem -= framerate
	}

	tm := t.Format("15:04:05")

	frames := strconv.Itoa(rem)
	if len(frames) < 2 {
		frames += "0"
	}

	tm += ":" + frames

	return tm, nil
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
