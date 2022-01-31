package main

import (
	"fmt"
	"hls"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error
type tickMsg time.Time

type bar struct {
	progress.Model
	name    string
	result  chan hls.Result
	len     int
	idx     int
	url     string
	outfile string
}
type model struct {
	textInput  *input
	bars       []*bar
	currentBar *bar
	err        error
}

func initialModel() model {
	ti := newInput()
	ti.Focus()
	return model{
		textInput: ti,
		err:       nil,
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			switch m.textInput.state {

			case insertURL:
				m.currentBar = &bar{url: m.textInput.Value()}
				m.textInput.next()

			case insertFile:
				b := m.currentBar
				b.outfile = m.textInput.Value()
				b.name = b.outfile

				p := progress.New(progress.WithDefaultGradient())
				b.Model = p
				var err error
				b.len, b.result, err = hls.Download(m.currentBar.url, m.currentBar.outfile)
				if err != nil {
					b.name += fmt.Sprintf(" - FAILED: %v", err)
				}

				m.bars = append(m.bars, b)
				m.textInput.next()
			}
		}

	case tickMsg:
		for _, b := range m.bars {
			select {
			case res := <-b.result:
				if res.Err != nil {
					b.name += fmt.Sprintf(" - FAILED: %v", res.Err)
					continue
				}
				b.idx++
				b.SetPercent(float64(b.idx) / float64(b.len))
				if b.idx == b.len {
					b.name += " - COMPLETED"
				}
			default:
			}
		}
		return m, tickCmd()

	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput.Model, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := m.textInput.View()
	for _, b := range m.bars {
		if b == nil {
			continue
		}
		s = fmt.Sprintf(
			"%s\n\n%s 	%s",
			s,
			b.ViewAs(b.Percent()),
			b.name,
		)
	}
	return s
}
