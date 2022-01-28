package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("Hello world")

	// f, err := os.Open("v2.m3u")
	// if err != nil {
	// 	panic(err)
	// }
	// p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	// if err != nil {
	// 	panic(err)
	// }
	// switch listType {
	// case m3u8.MEDIA:
	// 	mediapl := p.(*m3u8.MediaPlaylist)
	// 	for _, s := range mediapl.Segments {
	// 		if s != nil {
	// 			fmt.Println(s.URI)
	// 		}
	// 	}
	// 	fmt.Println(mediapl.Key)
	// }

	// validate := func(input string) error {
	// 	_, err := strconv.ParseFloat(input, 64)
	// 	if err != nil {
	// 		return errors.New("Invalid number")
	// 	}
	// 	return nil
	// }

	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
