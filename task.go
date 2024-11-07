package main

import (
	"strings"
)

type Task struct {
	title       string
	description string
	status      status
	project     string
	id          int
	urgency     int
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
	if t.urgency > 0 {
		return strings.Repeat("✦", t.urgency) + "\n" + t.title
	}
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type TaskList []Task

func (tl TaskList) Less(i, j int) bool {
	return (tl[i].urgency < tl[j].urgency) || (tl[i].urgency == tl[j].urgency && tl[i].title < tl[j].title)
}

func (tl TaskList) Len() int {
	return len(tl)
}

func (tl TaskList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (t *Task) AddUrgency() {
	if t.urgency < 10 {
		t.urgency++
	}
}

func (t *Task) SubtractUrgency() {
	if t.urgency > 0 {
		t.urgency--
	}
}
