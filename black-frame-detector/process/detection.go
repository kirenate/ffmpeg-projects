package process

import (
	"github.com/pkg/errors"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
)

func GetVideoMetadata(filename string) (chan []byte, chan error, chan interface{}) {
	wChan := make(chan []byte)
	errChan := make(chan error)
	done := make(chan interface{})

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		errChan <- errors.Wrap(err, "pipe create")
		return nil, nil, nil
	}
	os.Stdout = w

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
		errChan <- errors.Wrap(err, "cmd failed")
		return nil, nil, nil
	}

	err = w.Close()
	if err != nil {
		errChan <- errors.Wrap(err, "close w")
		return nil, nil, nil
	}

	go func() {

		var errLoop error
		for !errors.Is(errLoop, io.EOF) {
			buf := make([]byte, 10)

			_, errLoop = io.ReadAtLeast(r, buf, 10)
			if errLoop != nil && !errors.Is(errLoop, io.EOF) {
				errChan <- errLoop
				break
			}

			wChan <- buf

		}

		os.Stdout = oldStdout
		close(wChan)
		close(errChan)
		close(done)
	}()

	return wChan, errChan, done

}

//
//func ExtractBlackFramesFromMetadata(meta []byte) ([]byte, error) {
//}
