package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"main.go/process"
	"main.go/tui"
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

	m := tui.NewModel()
	p := tea.NewProgram(m)

	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()

		_, err := p.Run()
		if err != nil {
			panic(err)
		}
	}()

	in, err := process.GetVideoMetadata(inputPath, &wg)
	if err != nil {
		panic(err)
	}

	out := process.ExtractBlackFramesFromMetadata(in, &wg)

	timecode := process.DurationToTimestamp(inputPath, out, &wg)

	for data := range timecode {
		fmt.Println("sent", data)
		p.Send(data)
	}

	wg.Wait()
}

// -i{inputPath}
