package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type status int

func FromInt(i int) status {
	switch i {
	case 0:
		return todo
	case 1:
		return inProgress
	case 2:
		return done
	default:
		panic("invalid status")
	}
}

func (s status) getNext() status {
	if s == done {
		return todo
	}
	return s + 1
}

func (s status) getPrev() status {
	if s == todo {
		return done
	}
	return s - 1
}

const margin = 4

var board *Board

const (
	todo status = iota
	inProgress
	done
)

var project = "default"

func main() {
	args := os.Args
	if len(args) == 2 {
		project = args[1]
	} else if len(args) > 2 {
		fmt.Println("Project name must be one word")
		os.Exit(1)
	}
	err := initDatabase("./project/tasks.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	board = NewBoard()
	board.initLists()
	p := tea.NewProgram(board)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}