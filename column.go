package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const APPEND = -1

type column struct {
	list      list.Model
	status    status
	height    int
	width     int
	focus     bool
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
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#00ffd7", Dark: "#00ffff"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#00ffd7", Dark: "#00ffff"}).
		Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#00ffd7", Dark: "#00ffff"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#00ffd7", Dark: "#5fafd7"}).
		Padding(0, 0, 0, 1)

	defaultList := list.New([]list.Item{}, delegate, 0, 0)
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
	// Capture Window Size, resize the column and the project list accordingly
	case tea.WindowSizeMsg:
		c.setSize(msg.Width)
		projects.SetSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/margin, msg.Height-8)
		// Capture Key
	case tea.KeyMsg:
		if !c.filtering {
			// If its filtering, that allows you to type the keys that are shortcuts
			switch {
			case key.Matches(msg, keys.Edit):
				// Handles editting task
				if len(c.list.VisibleItems()) != 0 {
					task := c.list.SelectedItem().(Task)
					f := NewForm()
					f.Fill(&task)
					f.index = c.list.Index()
					f.col = c
					return f.Update(nil)
				}
			case key.Matches(msg, keys.New):
				// Handles new task
				f := NewForm()
				f.index = APPEND
				f.col = c
				return f.Update(nil)
			case key.Matches(msg, keys.Delete):
				// Deletes the task
				return c, c.DeleteCurrent()
			case key.Matches(msg, keys.Enter):
				// Moves the task to the next column
				return c, c.Move(c.status.getNext())
			case key.Matches(msg, keys.Prev):
				// Moves the task to the previous column
				return c, c.Move(c.status.getPrev())
			case key.Matches(msg, keys.Todo):
				// Moves the task to the todo column
				return c, c.Move(todo)
			case key.Matches(msg, keys.InProgress):
				// Moves the task to the in progress column
				return c, c.Move(inProgress)
			case key.Matches(msg, keys.Done):
				// Moves the task to the done column
				return c, c.Move(done)
			case key.Matches(msg, keys.AddUrgency):
				c.addUrgency()
			case key.Matches(msg, keys.RemoveUrgency):
				c.removeUrgency()
			case key.Matches(msg, keys.Refresh):
				board.resetLists()
			}
		}
		switch {
		// toggle filtering
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

func (c *column) captureItem() (Task, bool) {
	// get the selected item
	selectedItem := c.list.SelectedItem()
	if selectedItem == nil {
		// if there is no selected item, return false
		return Task{}, false
	}
	// return the selected item as a Task
	return selectedItem.(Task), true
}

func (c *column) DeleteCurrent() tea.Cmd {
	if len(c.list.VisibleItems()) > 0 {
		selectedTask, ok := c.captureItem()
		if !ok {
			return nil
		}
		// delete the task from the database
		deleteTask(&selectedTask)
		c.list.RemoveItem(c.list.Index())
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)
	return cmd
}

func (c *column) Set(i int, t Task) tea.Cmd {
	if i != APPEND {
		// if the index is not -1, then we are editing
		old := c.list.Items()[i].(Task)
		// edit the task in the database
		editTask(&t, &old)
		return c.list.SetItem(i, t)
	}
	// add the task to the database
	addTask(&t)
	projects.ResetProjects()
	return c.list.InsertItem(APPEND, t)
}

func (c *column) setSize(width int) {
	c.width = width / margin
}

func (c *column) getStyle() lipgloss.Style {
	if c.Focused() {
		// if the column is focused, itll have a different style
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	// if the column is not focused, itll be dimmer
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

func (c *column) Move(s status) tea.Cmd {
	// get the selected task
	selectedTask, ok := c.captureItem()
	if !ok {
		return nil
	}
	// get the current status of the task
	currentStatus := selectedTask.status
	// if the current status is the same as the new status, do nothing
	if currentStatus == s {
		return nil
	}
	// remove the task from the current, and add it to the new status
	c.list.RemoveItem(c.list.Index())
	selectedTask.status = s
	updateTaskStatus(&selectedTask)
	board.cols[s].list.
		InsertItem(
			len(c.list.Items())-1,
			list.Item(selectedTask),
		)
	return nil
}

func (c *column) addUrgency() {
	// get the selected task
	selectedTask, ok := c.captureItem()
	if !ok {
		return
	}
	// add urgency to the task
	selectedTask.AddUrgency()
	// update the task in the database
	changeUrgency(&selectedTask)
	// update the task in the list
	c.list.SetItem(c.list.Index(), selectedTask)
}

func (c *column) removeUrgency() {
	// get the selected task
	selectedTask, ok := c.captureItem()
	if !ok {
		return
	}
	// add urgency to the task
	selectedTask.SubtractUrgency()
	// update the task in the database
	changeUrgency(&selectedTask)
	// update the task in the list
	c.list.SetItem(c.list.Index(), selectedTask)
}
