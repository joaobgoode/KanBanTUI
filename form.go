package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	taskTitle = iota
	taskDate
	taskShortDesc
	taskLongDesc
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff"))
	continueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
)

func titleValidador(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("Title is required")
	}
	return nil
}

func dateValidador(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("Date is required")
	}
	if len(s) != 10 {
		return fmt.Errorf("Date must be in format DD/MM/YYYY")
	}
	if s[2] != '/' || s[5] != '/' {
		return fmt.Errorf("Date must be in format DD/MM/YYYY")
	}
	numbers := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(numbers, 10, 64)
	if err != nil {
		return fmt.Errorf("Date must be in format DD/MM/YYYY")
	}
	return nil
}

type Form struct {
	err     error
	inputs  []textinput.Model
	col     column
	focused int
	index   int
}

func NewForm() *Form { // Updated to return *Form
	inputs := make([]textinput.Model, 4)
	inputs[taskTitle] = textinput.New()
	inputs[taskTitle].Placeholder = "task title"
	inputs[taskTitle].Focus()
	inputs[taskTitle].CharLimit = 20
	inputs[taskTitle].Width = 30
	inputs[taskTitle].Prompt = ""
	inputs[taskTitle].Validate = titleValidador

	inputs[taskDate] = textinput.New() // Ensure each input is initialized
	inputs[taskDate].Placeholder = "DD/MM/YYYY"
	inputs[taskDate].CharLimit = 10
	inputs[taskDate].Width = 30
	inputs[taskDate].Prompt = ""
	inputs[taskDate].Validate = dateValidador

	inputs[taskShortDesc] = textinput.New()
	inputs[taskShortDesc].Placeholder = "Short description"
	inputs[taskShortDesc].CharLimit = 20
	inputs[taskShortDesc].Width = 30
	inputs[taskShortDesc].Prompt = ""

	inputs[taskLongDesc] = textinput.New()
	inputs[taskLongDesc].Placeholder = "Long description"
	inputs[taskLongDesc].CharLimit = 100
	inputs[taskLongDesc].Width = 30
	inputs[taskLongDesc].Prompt = ""

	return &Form{inputs: inputs, focused: 0, err: nil}
}

func (f Form) CreateTask() Task {
	t := Task{
		status:      f.col.status,
		title:       f.inputs[taskTitle].Value(),
		date:        f.inputs[taskDate].Value(),
		description: f.inputs[taskShortDesc].Value(),
		longdesc:    f.inputs[taskShortDesc].Value(),
		urgency:     0,
		project:     project,
	}
	log.Println(t)
	return t
}

func (f *Form) Fill(task *Task) {
	log.Println("Fill")
	f.inputs[taskTitle].SetValue(task.title)
	log.Println("title")
	f.inputs[taskDate].SetValue(task.getDate())
	log.Println("date")
	f.inputs[taskShortDesc].SetValue(task.description)
	log.Println("desc")
	f.inputs[taskLongDesc].SetValue(task.longdesc)
	log.Println("long")
}

func (f *Form) nextInput() {
	log.Println(len(f.inputs))
	f.focused = (f.focused + 1) % len(f.inputs)
}

func (f *Form) prevInput() { // Changed to pointer receiver
	f.focused--
	// Wrap around
	if f.focused < 0 {
		f.focused = len(f.inputs) - 1
	}
}

func (f Form) Init() tea.Cmd { // Changed to pointer receiver
	return textinput.Blink
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) { // Changed to pointer receiver
	cmds := make([]tea.Cmd, len(f.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			// Check for validation errors when hitting Enter on the last field
			if f.focused == taskLongDesc {
				err := f.inputs[taskTitle].Validate(f.inputs[taskTitle].Value())
				if err != nil {
					f.err = err
					return f, nil
				}
				err = f.inputs[taskDate].Validate(f.inputs[taskDate].Value())
				if err != nil {
					f.err = err
					return f, nil
				}
				return board.Update(f)
			}
			f.nextInput()
		case key.Matches(msg, keys.Up):
			f.prevInput()
		case key.Matches(msg, keys.Down):
			f.nextInput()
		case key.Matches(msg, keys.Back):
			// if the back key is pressed, return to the board
			return board.Update(nil)

		}

		for i := range f.inputs {
			f.inputs[i].Blur()
		}
		f.inputs[f.focused].Focus()
	}

	for i := range f.inputs {
		f.inputs[i], cmds[i] = f.inputs[i].Update(msg)
	}
	return f, tea.Batch(cmds...)
}

func (f Form) View() string { // Changed to pointer receiver
	errorMessage := ""
	if f.err != nil {
		errorMessage = errorStyle.Render(f.err.Error())
	}

	return fmt.Sprintf(
		` New task:

 %s
 %s

 %s
 %s

 %s
 %s
 
 %s
 %s
 
 %s
 
 %s
`,
		inputStyle.Width(30).Render("Task Title"),
		f.inputs[taskTitle].View(),
		inputStyle.Width(30).Render("Date:"),
		f.inputs[taskDate].View(),
		inputStyle.Width(30).Render("Short Description"),
		f.inputs[taskShortDesc].View(),
		inputStyle.Width(30).Render("Long Description"),
		f.inputs[taskLongDesc].View(),
		continueStyle.Render("Continue ->"),
		errorMessage, // Render error message in red
	) + "\n"
}
