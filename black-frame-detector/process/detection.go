package process

import (
	"github.com/pkg/errors"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
)

func GetVideoMetadata(filename string, wg *sync.WaitGroup) (chan []byte, error) {
	logger := log.Default()
	logger.SetOutput(os.Stderr)

	out := make(chan []byte)

	r, w, err := os.Pipe()
	if err != nil {
		return nil, errors.Wrap(err, "pipe create")
	}

	inputKwargs := ffmpeg.KwArgs{}
	inputKwargs["hide_banner"] = ""

	outputKwargs := ffmpeg.KwArgs{}
	outputKwargs["filter:v"] = "blackframe"
	outputKwargs["f"] = "null"

	cmd := ffmpeg.Input(filename, inputKwargs).
		Output("out.null", outputKwargs).
		WithErrorOutput(w)

	err = cmd.Run()
	if err != nil {
		return nil, errors.Wrap(err, "cmd failed")
	}

	err = w.Close()
	if err != nil {
		return nil, errors.Wrap(err, "close w")
	}

	go func() {
		defer wg.Done()

		for !errors.Is(err, io.ErrUnexpectedEOF) {
			buf := make([]byte, 10)

			_, err = io.ReadAtLeast(r, buf, 10)
			if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
				logger.Println("readAtLeast: ", err.Error())
				break
			}

			out <- buf

		}

		close(out)
	}()

	return out, nil

}

func ExtractBlackFramesFromMetadata(in <-chan []byte, wg *sync.WaitGroup) chan string {
	out := make(chan string)

	reg := regexp.MustCompile(` frame:[0-9]{3,}`)

	go func() {
		defer wg.Done()

		var str []byte
		var tmp string
		for data := range in {
			str = slices.Concat(str, data)
			stamp := reg.Find(str)
			if stamp != nil {
				tmp, _ = strings.CutPrefix(string(stamp), " frame:")
				out <- tmp
				str = []byte{}
			}
		}

		close(out)
	}()

	return out
}
