package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDatabase(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			title TEXT NOT NULL, 
			description TEXT, 
      project TEXT NOT NULL,
			status REAL NOT NULL
		)`,
	)
	if err != nil {
		return err
	}
	return nil
}

func addTask(t *Task) {
	s := int(t.status)
	res, err := db.ExecContext(
		context.Background(),
		`INSERT INTO tasks (title, description, project, status) VALUES (?,?,?,?);`, t.title, t.description, project, s,
	)
	if err != nil {
		panic(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	t.id = int(id)
	// rest of the function
}

func taskByStatus(status status) ([]Task, error) {
	// this slice will contain all the albums retrieved

	var tasks []Task
	value := int(status)
	rows, err := db.QueryContext(
		context.Background(),
		`SELECT * FROM tasks WHERE status=? AND project=?;`, value, project,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// rest of the function
	for rows.Next() {
		var t Task
		var s int
		err := rows.Scan(&t.id, &t.title, &t.description, &t.project, &s)
		if err != nil {
			return nil, err
		}
		t.status = FromInt(s)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func updateTaskStatus(t *Task) {
	_, err := db.ExecContext(
		context.Background(),
		`UPDATE tasks SET status=? WHERE  id = ?;`, t.status, t.id,
	)
	if err != nil {
		panic(err)
	}
}

func deleteTask(t *Task) {
	_, err := db.ExecContext(
		context.Background(),
		`DELETE FROM tasks WHERE id=?;`, t.id,
	)
	if err != nil {
		panic(err)
	}
}

func editTask(t *Task, old *Task) {
	_, err := db.ExecContext(
		context.Background(),
		`UPDATE tasks SET title=?, description=? WHERE id=?;`, t.title, t.description, old.id,
	)
	if err != nil {
		panic(err)
	}
}

func getProjects() ([]projectItem, error) {
	rows, err := db.QueryContext(
		context.Background(),
		`
	SELECT project, COUNT(*) AS project_count
	FROM tasks
	GROUP BY project;
	`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	projects := []projectItem{}
	for rows.Next() {
		var title string
		var count int
		err := rows.Scan(&title, &count)
		if err != nil {
			return nil, err
		}
		description := fmt.Sprintf("%d tasks", count)
		projects = append(projects, projectItem{title, description})
	}
	return projects, nil
}

func checkProject(project string) bool {
	rows, err := db.QueryContext(
		context.Background(),
		`SELECT * FROM tasks WHERE project=?;`, project,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	return rows.Next()
}
