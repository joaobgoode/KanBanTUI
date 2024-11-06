package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	help        help.Model
	title       textinput.Model
	description textarea.Model
	col         column
	index       int
}

func newDefaultForm() *Form {
	return NewForm("task name", "", false)
}

func NewForm(title, description string, edit bool) *Form {
	form := Form{
		help:        help.New(),
		title:       textinput.New(),
		description: textarea.New(),
	}
	form.title.Placeholder = title
	form.description.Placeholder = description
	form.title.Focus()
	if edit {
		// if eddit is true, preset the title and description
		form.title.SetValue(title)
		form.description.SetValue(description)
	}
	return &form
}

func (f Form) CreateTask() Task {
	t := Task{
		status:      f.col.status,
		title:       f.title.Value(),
		description: f.description.Value(),
		project:     project,
	}
	return t
}

func (f Form) Init() tea.Cmd {
	return nil
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case column:
		f.col = msg
		f.col.list.Index()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			// if the back key is pressed, return to the board
			return board.Update(nil)
		case key.Matches(msg, keys.Enter):
			if f.title.Focused() {
				f.title.Blur()
				f.description.Focus()
				return f, textarea.Blink
			}
			// Return the completed form as a message.
			return board.Update(f)
		}
	}
	if f.title.Focused() {
		f.title, cmd = f.title.Update(msg)
		return f, cmd
	}
	f.description, cmd = f.description.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Create a new task",
		f.title.View(),
		f.description.View(),
		f.help.View(keys))
}
