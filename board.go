package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	cols      []column
	focused   status
	loaded    bool
	quitting  bool
	filtering bool
}

func NewBoard() *Board {
	return &Board{focused: todo}
}

func (m *Board) Init() tea.Cmd {
	return nil
}

func (m *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		for i := 0; i < len(m.cols); i++ {
			var res tea.Model
			res, cmd = m.cols[i].Update(msg)
			m.cols[i] = res.(column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		if project == "" {
			return projects.Update(nil)
		}
		return m, tea.Batch(cmds...)
	case Form:
		return m, m.cols[m.focused].Set(msg.index, msg.CreateTask())
	case moveMsg:
		return m, m.cols[m.focused.getNext()].Set(APPEND, msg.Task)
	case tea.KeyMsg:
		if !m.filtering {
			switch {
			case key.Matches(msg, keys.Quit):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, keys.Left):
				m.cols[m.focused].Blur()
				m.focused = m.focused.getPrev()
				m.cols[m.focused].Focus()
			case key.Matches(msg, keys.Right):
				m.cols[m.focused].Blur()
				m.focused = m.focused.getNext()
				m.cols[m.focused].Focus()
			case key.Matches(msg, keys.Projects):
				return projects.Update(msg)
			}
		}
		switch {
		case key.Matches(msg, keys.Filtering):
			m.filtering = true
		case key.Matches(msg, keys.Back), key.Matches(msg, keys.Enter):
			m.filtering = false
		}
	}
	res, cmd := m.cols[m.focused].Update(msg)
	if _, ok := res.(column); ok {
		m.cols[m.focused] = res.(column)
	} else {
		return res, cmd
	}
	return m, cmd
}

var cheatsheetStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Foreground(lipgloss.Color("#767676"))

// Changing to pointer receiver to get back to this model after adding a new task via the form... Otherwise I would need to pass this model along to the form and it becomes highly coupled to the other models.
func (m *Board) View() string {
	// clears the screen before rendering
	if m.quitting {
		return ""
	}
	// if the board is not loaded, return loading
	if !m.loaded {
		return "loading..."
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[todo].View(),
		m.cols[inProgress].View(),
		m.cols[done].View(),
		cheatsheetStyle.Render(
			"Movement:\n"+
				"←/→: switch columns\n"+
				"↑/↓ or j/k: move up and down\n"+
				"\nTasks:\n"+
				"enter: move task forward\n"+
				"b: move task backward\n"+
				"1: move task to To Do\n"+
				"2: move task to In Progress\n"+
				"3: move task to Done\n"+
				"+/-: add/remove urgency\n"+
				"n: new task\n"+
				"e: edit task\n"+
				"del: delete task\n"+
				"\nFilter:\n"+
				"enter: leave with filter on\n"+
				"esc: leave with filter off\n"+
				"\nProjects:\n"+
				"p: see projects",
		),
	)
	return lipgloss.JoinVertical(lipgloss.Left, board)
}
