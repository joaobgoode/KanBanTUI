package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const APPEND = -1

type column struct {
	focus     bool
	status    status
	list      list.Model
	height    int
	width     int
	filtering bool
}

func (c *column) Focus() {
	c.focus = true
}

func (c *column) Blur() {
	c.focus = false
}

func (c *column) Focused() bool {
	return c.focus
}

func newColumn(status status) column {
	var focus bool
	if status == todo {
		focus = true
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return column{focus: focus, status: status, list: defaultList}
}

// Init does initial setup for the column.
func (c column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		projects.SetSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/margin, msg.Height-8)
	case tea.KeyMsg:
		if !c.filtering {
			switch {
			case key.Matches(msg, keys.Edit):
				if len(c.list.VisibleItems()) != 0 {
					task := c.list.SelectedItem().(Task)
					f := NewForm(task.title, task.description)
					f.index = c.list.Index()
					f.col = c
					return f.Update(nil)
				}
			case key.Matches(msg, keys.New):
				f := newDefaultForm()
				f.index = APPEND
				f.col = c
				return f.Update(nil)
			case key.Matches(msg, keys.Delete):
				return c, c.DeleteCurrent()
			case key.Matches(msg, keys.Enter):
				return c, c.MoveToNext()
			case key.Matches(msg, keys.Prev):
				return c, c.MoveToPrev()
			case key.Matches(msg, keys.Todo):
				return c, c.MoveToTodo()
			case key.Matches(msg, keys.InProgress):
				return c, c.MoveToInProgress()
			case key.Matches(msg, keys.Done):
				return c, c.MoveToDone()
			}
		}
		switch {
		case key.Matches(msg, keys.Filtering):
			c.filtering = true
		case key.Matches(msg, keys.Back), key.Matches(msg, keys.Enter):
			c.filtering = false
		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c column) View() string {
	return c.getStyle().Render(c.list.View())
}

func (c *column) DeleteCurrent() tea.Cmd {
	if len(c.list.VisibleItems()) > 0 {
		selectedItem := c.list.SelectedItem()
		if selectedItem == nil {
			return nil
		}
		selectedTask := selectedItem.(Task)
		deleteTask(&selectedTask)
		c.list.RemoveItem(c.list.Index())
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)
	return cmd
}

func (c *column) Set(i int, t Task) tea.Cmd {
	if i != APPEND {
		old := c.list.Items()[i].(Task)
		editTask(&t, &old)
		return c.list.SetItem(i, t)
	}
	addTask(&t)
	projects.ResetProjects()
	return c.list.InsertItem(APPEND, t)
}

func (c *column) setSize(width, height int) {
	c.width = width / margin
}

func (c *column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#474a51")).
		Height(c.height).
		Width(c.width)
}

type moveMsg struct {
	Task
}

func (c *column) MoveToPrev() tea.Cmd {
	selectedItem := c.list.SelectedItem()
	if selectedItem == nil {
		return nil
	}
	selectedTask := selectedItem.(Task)
	if selectedTask.status == todo {
		return nil
	}
	c.list.RemoveItem(c.list.Index())
	selectedTask.status = selectedTask.status.getPrev()
	updateTaskStatus(&selectedTask)
	board.cols[selectedTask.status].list.
		InsertItem(
			len(c.list.Items())-1,
			list.Item(selectedTask),
		)
	return nil
}

func (c *column) MoveToNext() tea.Cmd {
	selectedItem := c.list.SelectedItem()
	if selectedItem == nil {
		return nil
	}
	selectedTask := selectedItem.(Task)
	if selectedTask.status == done {
		return nil
	}
	c.list.RemoveItem(c.list.Index())
	selectedTask.status = selectedTask.status.getNext()
	updateTaskStatus(&selectedTask)
	board.cols[selectedTask.status].list.
		InsertItem(
			len(c.list.Items())-1,
			list.Item(selectedTask),
		)
	return nil
}

func (c *column) MoveToTodo() tea.Cmd {
	selectedItem := c.list.SelectedItem()
	if selectedItem == nil {
		return nil
	}
	selectedTask := selectedItem.(Task)
	if selectedTask.status == todo {
		return nil
	}
	c.list.RemoveItem(c.list.Index())
	selectedTask.status = todo
	updateTaskStatus(&selectedTask)
	board.cols[todo].list.
		InsertItem(
			len(c.list.Items())-1,
			list.Item(selectedTask),
		)
	return nil
}

func (c *column) MoveToInProgress() tea.Cmd {
	selectedItem := c.list.SelectedItem()
	if selectedItem == nil {
		return nil
	}
	selectedTask := selectedItem.(Task)
	if selectedTask.status == inProgress {
		return nil
	}
	c.list.RemoveItem(c.list.Index())
	selectedTask.status = inProgress
	updateTaskStatus(&selectedTask)
	board.cols[inProgress].list.
		InsertItem(
			len(c.list.Items())-1,
			list.Item(selectedTask),
		)
	return nil
}

func (c *column) MoveToDone() tea.Cmd {
	selectedItem := c.list.SelectedItem()
	if selectedItem == nil {
		return nil
	}
	selectedTask := selectedItem.(Task)
	if selectedTask.status == done {
		return nil
	}
	c.list.RemoveItem(c.list.Index())
	selectedTask.status = done
	updateTaskStatus(&selectedTask)
	board.cols[done].list.
		InsertItem(
			len(c.list.Items())-1,
			list.Item(selectedTask),
		)
	return nil
}
