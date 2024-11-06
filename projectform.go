package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type projectForm struct {
	help     help.Model
	title    textinput.Model
	valid    bool
	editting bool
}

func newDefaultProjectForm() *projectForm {
	return NewProjectForm("project name", false)
}

func NewProjectForm(title string, edit bool) *projectForm {
	form := projectForm{
		help:     help.New(),
		title:    textinput.New(),
		valid:    true,
		editting: edit,
	}
	form.title.Placeholder = title
	if edit {
		form.title.SetValue(form.title.Placeholder)
	}
	form.title.CharLimit = 30
	form.title.Focus()
	return &form
}

func (f projectForm) Init() tea.Cmd {
	return nil
}

func (f projectForm) ValidTitle() bool {
	return f.title.Value() != "" && !strings.ContainsAny(f.title.Value(), " \t\n\r") && !checkProject(f.title.Value())
}

func (f projectForm) ValidEdit() bool {
	return f.title.Value() != "" && !strings.ContainsAny(f.title.Value(), " \t\n\r") && (!checkProject(f.title.Value()) || f.title.Value() == f.title.Placeholder)
}

func (f projectForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			return projects.Update(nil)
		case key.Matches(msg, keys.Enter):
			if !f.editting && f.ValidTitle() {
				f.valid = true
				project = f.title.Value()
				board.resetLists()
				return board.Update(nil)
			} else if f.editting && f.ValidEdit() {
				f.valid = true
				changeProject(f.title.Placeholder, f.title.Value())
				board.resetLists()
				projects.ResetProjects()
				return projects.Update(nil)
			} else {
				f.valid = false
			}
		}
	}
	f.title, cmd = f.title.Update(msg)
	return f, cmd
}

func (f projectForm) getStyle() lipgloss.Style {
	if f.valid {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("2"))
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("9"))
}

func (f projectForm) View() string {
	return f.getStyle().Render(f.title.View())
}
