package main

import "github.com/charmbracelet/bubbles/textinput"

var labels []string = []string{
	"Insert URL of HLS playlist to download",
	"Insert destination filepath",
}

type state int

const (
	insertURL state = iota
	insertFile
	none
)

type input struct {
	textinput.Model
	state state
	label string
}

func newInput() *input {
	i := &input{Model: textinput.New()}
	i.state = insertURL
	i.label = labels[i.state]
	i.Placeholder = i.label
	return i
}

func (i *input) next() {
	if i.state++; i.state == none {
		i.state = insertURL
	}
	i.label = labels[i.state]
	i.Placeholder = i.label
	i.Reset()
}
