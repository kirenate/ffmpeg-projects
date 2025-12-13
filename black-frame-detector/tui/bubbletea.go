package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"
	"os"
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

func mustLog(v string) {
	file, err := os.OpenFile("/tmp/out.log", os.O_RDWR, 0)
	if err != nil {
		file, err = os.Create("/tmp/out.log")
		if err != nil {
			panic(errors.Wrap(err, "failed to open and create"))
		}
	}

	defer file.Close()

	_, err = file.WriteString(v)
	if err != nil {
		panic(errors.Wrap(err, "failed to write string"))
	}
}

func (m Model) View() string {
	mustLog(fmt.Sprintf("%+v", m.choices))

	s := "Detecting black frames. . .\n\n"
	var prefix string

	for i, data := range m.choices {
		prefix = "   "
		if _, ok := m.selected[i]; ok {
			prefix = "---"
		}

		s += fmt.Sprintf("%v %s %s\n", i, prefix, data)
	}

	s += "\nPress ctrl+C or q to exit\n"

	mustLog(s)
	return s
}
