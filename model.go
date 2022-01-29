package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type bar struct {
	progress.Model
	name string
	prog chan struct{}
	len  int
}

type errMsg error

type model struct {
	textInput textinput.Model
	bars      []bar
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	// return textinput.Blink
	return tickCmd()
}

var url string
var file string
var inpstat = 1

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			switch inpstat {
			case 1:
				// url = m.textInput.Value()
				// m.textInput.Reset()
				// inpstat = 2
			case 2:
				// file = m.textInput.Value()
				// c, l := download(url, file, 10)
				// p := progress.New(progress.WithDefaultGradient(), progress.WithWidth(l))
				// m.bars = append(m.bars, bar{Model: p, name: m.textInput.Value(), prog: c})
				// m.textInput.Reset()
				// inpstat = 1
			}
		}

		// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		// progressModel, cmd := m.progress.Update(msg)
		// m.progress = progressModel.(progress.Model)
		// return m, cmd

	case tickMsg:
		var cmd tea.Cmd
		for _, p := range m.bars {
			select {
			case <-p.prog:
				fmt.Println("inc")
				cmd = p.IncrPercent(0.1)
			default:
			}
		}
		return m, tea.Batch(tickCmd(), cmd)

	case tea.WindowSizeMsg:
		for _, p := range m.bars {
			p.Width = msg.Width / 2
			// if p.Width > maxWidth {
			// 	p.Width = maxWidth
			// }
		}
		return m, nil

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// return fmt.Sprintf(
	// 	"What’s your favorite Pokémon?\n\n%s\n\n%s",
	// 	m.textInput.View(),
	// 	"(esc to quit)",
	// ) + "\n"
	s := m.textInput.View()
	for _, b := range m.bars {
		s = fmt.Sprintf(
			"%s\n%s 	%s",
			s,
			b.View(),
			b.name,
		)
	}
	return s
}
