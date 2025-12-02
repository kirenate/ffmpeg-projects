package tests_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"main.go/process"
	"sync"
	"testing"
	"time"
)

func TestGetVideoTimestamp(t *testing.T) {
	info, err := process.GetVideoInfo(Filename)
	require.NoError(t, err)
	fmt.Println(info)
	require.NotEmpty(t, info)
}

func TestGetBlackFramesTimestamps(t *testing.T) {
	var wg sync.WaitGroup
	var str string

	wg.Add(1)
	wChan, err := process.GetVideoMetadata(Filename, &wg)
	require.NoError(t, err)

	str = checker(str, wChan)
	wg.Wait()
	require.NotEmpty(t, str)

	fmt.Println(err)
	fmt.Println("str:", str)
}

func checker(str string, wChan chan []byte) string {
	fmt.Println("1")
	for data := range wChan {
		str += string(data)
	}

	return str
}

func TestExtractBlackFramesFromMetadata(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(2)
	in, err := process.GetVideoMetadata(Filename, &wg)
	require.NoError(t, err)

	out := process.ExtractBlackFramesFromMetadata(in, &wg)

	var str []string
	for range out {
		str = append(str, <-out)
	}

	wg.Wait()

	require.NotEmpty(t, str)

	fmt.Println(str)
	time.Sleep(3 * time.Second)
}

func TestDurationToTimestamp(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(3)
	in, err := process.GetVideoMetadata(Filename, &wg)
	require.NoError(t, err)

	out := process.ExtractBlackFramesFromMetadata(in, &wg)

	timecode := process.DurationToTimestamp(Filename, out, &wg)

	var res []string

	for data := range timecode {
		res = append(res, data)
	}

	wg.Wait()

	require.NotEmpty(t, res)
	fmt.Println(res)

	time.Sleep(3 * time.Second)
}

var Filename = "./../black_filter3.mp4"
