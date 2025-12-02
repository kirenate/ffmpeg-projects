package main

import (
	"fmt"
	"main.go/process"
	"os"
	"sync"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		panic("Usage: bfd -i {inputPath}\n\nThis script only works with mp4 and mov")
	}
	if args[1] != "-i" {
		panic("Usage: bfd -i {inputPath}\n\nThis script only works with mp4 and mov")
	}

	inputPath := args[2]

	var wg sync.WaitGroup

	wg.Add(3)
	in, err := process.GetVideoMetadata(inputPath, &wg)
	if err != nil {
		panic(err)
	}
	out := process.ExtractBlackFramesFromMetadata(in, &wg)

	timecode := process.DurationToTimestamp(inputPath, out, &wg)

	for data := range timecode {
		fmt.Println(data)
	}

	wg.Wait()

}

// -i{inputPath}
