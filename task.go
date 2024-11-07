package main

import (
	"fmt"
	"strings"
)

type Task struct {
	title       string
	date        string
	description string
	longdesc    string
	project     string
	status      status
	id          int
	urgency     int
}

func (t *Task) getDate() string {
	parts := strings.Split(t.date, "/")
	return fmt.Sprintf("%s/%s/%s", parts[2], parts[1], parts[0])
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
	return t.getDate() + "  " + t.description
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
