package main

type Task struct {
	title       string
	description string
	status      status
	project     string
	id          int
}

func NewTask(status status, title, description string) Task {
	t := Task{status: status, title: title, description: description}
	addTask(&t)
	return t
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}
