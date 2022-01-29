package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Hello world")

	// p := tea.NewProgram(initialModel())

	// if err := p.Start(); err != nil {
	// 	log.Fatal(err)
	// }

	var (
		dir  = "./test"
		name = "test"
	)

	b, err := download("http://d2zihajmogu5jn.cloudfront.net/bipbop-advanced/gear3/prog_index.m3u8")
	if err != nil {
		fmt.Printf("Error while downloading video: %v", err)
		os.Exit(1)
	}

	media, err := decode(b)
	if err != nil {
		fmt.Printf("Error while decoding video: %v", err)
		os.Exit(1)
	}

	for _, s := range media.Segments {
		// id := s.SeqId
		// fmt.Printf(filepath.Join(dir, fmt.Sprint(s.SeqId)))
		// return
		err := os.MkdirAll(dir, os.Mode)
		if err != nil {
			fmt.Printf("Error while creating directory to store media segments", err)
			// continue
			return
		}

		f, err := os.Create(filepath.Join(dir, fmt.Sprint(s.SeqId)))
		if err != nil {
			fmt.Printf("Error while creating file for segment %d of %s: %v", s.SeqId, name, err)
			// continue
			return
		}

		b, err := download(s.URI)
		if err != nil {
			fmt.Printf("Error while downloading segment %d of %s", s.SeqId, name)
		}

		_, err = f.Write(b)
		if err != nil {
			fmt.Printf("Error while writing segment %d of %s", s.SeqId, name)
		}
	}
}
