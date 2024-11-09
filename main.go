package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type status int

func FromInt(i int) status {
	if i >= 3 && i < 0 {
		panic("invalid status")
	}
	return status(i)
}

func (s status) getNext() status {
	// If the status is done, cannot go further
	if s == done {
		return done
	}
	return s + 1
}

func (s status) getPrev() status {
	// If the status is todo, cannot go backards
	if s == todo {
		return todo
	}
	return s - 1
}

const margin = 4

var (
	board    *Board
	projects *projectList
)

const (
	todo status = iota
	inProgress
	done
)

var project = ""

func main() {
	// handles logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	// get executable path
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	// check if another instance is running, if so exit
	// otherwise could lead to funky behavior due to
	// detabase being out of sync
	pidFile := filepath.Join(exPath, "tasks.pid")
	if _, err := os.Stat(pidFile); err == nil {
		fmt.Println("Another instance is already running")
		log.Fatal("Another instance is already running")
	}
	pid := strconv.Itoa(os.Getpid())
	err = os.WriteFile(pidFile, []byte(pid), 0644)
	if err != nil {
		log.Fatalf("Failed to create PID file: %v", err)
	}
	defer os.Remove(pidFile)

	// handles starting the program
	args := os.Args
	if len(args) == 2 {
		// if there is args, get the project with the args name
		project = args[1]
	} else if len(args) > 2 {
		// cannot have more than one word for project name
		fmt.Println("Project name must be one word")
		os.Exit(1)
	}
	// start the db
	dbPath := filepath.Join(exPath, "tasks.db")
	err = initDatabase(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// start the board and projects
	board = NewBoard()
	board.initLists()
	projects = NewProjectList()
	projects.LoadProjects()
	p := tea.NewProgram(board, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
