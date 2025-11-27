package tests_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"main.go/process"
	"testing"
)

func TestGetVideoTimestamp(t *testing.T) {
	info, err := process.GetVideoInfo("./../testdata.mp4")
	require.NoError(t, err)
	fmt.Println(info)
	require.NotEmpty(t, info)
}

func TestDurationToTimestamp(t *testing.T) {
	info, err := process.GetVideoInfo("./../output_red.mp4")
	require.NoError(t, err)
	require.NotEmpty(t, info)

	tm, err := process.DurationToTimestamp(info)
	require.NoError(t, err)
	require.NotEmpty(t, tm)

	fmt.Println(tm)
}
