package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func NewModel() Model {
	return Model{
		choices:  []string{},
		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	case string:
		m.choices = append(m.choices, msg)
	}
	return m, nil
}

func (m Model) View() string {
	var s string
	var prefix string

	for i, data := range m.choices {
		prefix = "   "
		if _, ok := m.selected[i]; ok {
			prefix = "---"
		}

		s += fmt.Sprintf("%s%s\n", prefix, data)
	}

	return s
}
