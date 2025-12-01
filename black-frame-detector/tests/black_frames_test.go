package tests_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"main.go/process"
	"testing"
	"time"
)

func TestGetVideoTimestamp(t *testing.T) {
	info, err := process.GetVideoInfo("./../testdata.mp4")
	require.NoError(t, err)
	fmt.Println(info)
	require.NotEmpty(t, info)
}

func TestDurationToTimestamp(t *testing.T) {
	info, err := process.GetVideoInfo("./../output_black.mp4")
	require.NoError(t, err)
	require.NotEmpty(t, info)

	tm, err := process.DurationToTimestamp(info)
	require.NoError(t, err)
	require.NotEmpty(t, tm)

	fmt.Println(tm)
}

func TestGetBlackFramesTimestamps(t *testing.T) {
	var str string

	wChan, errChan, done := process.GetVideoMetadata("./../black_filter3.mp4")
	time.Sleep(5 * time.Second)
	str, err := checker(str, wChan, errChan, done)
	time.Sleep(5 * time.Second)
	require.NotEmpty(t, str)
	fmt.Println(err)
	fmt.Println("str:", str)
}

func checker(str string, wChan chan []byte, errChan chan error, done chan interface{}) (string, error) {
	fmt.Println("1")
	for {
		select {
		case data := <-wChan:
			str += string(data)
		case err := <-errChan:
			time.Sleep(5 * time.Second)
			return str, err
		case <-done:
			time.Sleep(5 * time.Second)
			fmt.Println("done! returning")
			return str, nil
		}
	}

}
