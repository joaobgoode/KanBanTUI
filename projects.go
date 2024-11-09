package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type projectItem struct {
	title       string
	description string
}

// implement the list.Item interface
func (p projectItem) FilterValue() string {
	return p.title
}

func (p projectItem) Title() string {
	return p.title
}

func (p projectItem) Description() string {
	return p.description
}

type projectList struct {
	list      list.Model
	height    int
	width     int
	filtering bool
}

func (p *projectList) SetSize(width, height int) {
	p.width = width / margin
	p.height = height - 2*margin
	p.list.SetSize(p.width, p.height)
}

func (p *projectList) Init() tea.Cmd {
	return nil
}

func (p *projectList) getStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#17fc03")).
		Foreground(lipgloss.Color("2")).
		Height(p.height).
		Width(p.width)
}

func (p *projectList) View() string {
	return p.getStyle().Render(p.list.View())
}

func NewProjectList() *projectList {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#17fc03"))
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	emptyList := list.New([]list.Item{}, delegate, 0, 0)
	return &projectList{list: emptyList, filtering: false}
}

func (p *projectList) RemoveAllItems() {
	p.list.SetItems([]list.Item{})
}

func (p *projectList) LoadProjects() {
	projects, err := getProjects()
	if err != nil {
		panic(err)
	}
	for _, pi := range projects {
		p.list.InsertItem(-1, pi)
	}
	p.list.Title = "Projects"
}

func (p *projectList) ResetProjects() {
	p.RemoveAllItems()
	p.LoadProjects()
}

func (p *projectList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !p.filtering {
			switch {
			case key.Matches(msg, keys.New):
				pf := NewProjectForm("project name", false)
				return pf.Update(nil)
			case key.Matches(msg, keys.Edit):
				if len(p.list.VisibleItems()) != 0 {
					pi := p.list.SelectedItem().(projectItem)
					pf := NewProjectForm(pi.title, true)
					return pf.Update(nil)
				}
			case key.Matches(msg, keys.Quit):
				if project == "" {
					return p, tea.Quit
				}
				return board.Update(nil)
			}
		}
		switch {
		case key.Matches(msg, keys.Filtering):
			p.filtering = true
		case key.Matches(msg, keys.Back):
			p.filtering = false
		case key.Matches(msg, keys.Enter):
			p.filtering = false
			selected := p.list.SelectedItem()
			if selected == nil {
				return p.Update(nil)
			}
			project = selected.(projectItem).title
			board.resetLists()
			return board.Update(nil)
		}
	}
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}
