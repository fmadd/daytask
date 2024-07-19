package sqlite

import (
	"database/sql"
	"daytask/internal/storage"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS daytask(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		title TEXT NOT NULL,
		owner TEXT NOT NULL,
		date DATE NOT NULL);
		CREATE INDEX IF NOT EXIST idx_owner ON daytask(owner);
		`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	
	}
	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveTask(taskName string, taskOwner string, taskDate string) (int64, error){
	const op = "storage.sqlite.SaveTask"

	stmt, err := s.db.Prepare("INSERT INTO daytask(title, owner, date) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(taskName, taskOwner, taskDate)
	if err != nil { 
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil 
}

func (s *Storage) DeleteTask(taskName string, taskOwner string, taskDate string) (error){
	const op = "storage.sqlite.DeleteTask"

	stmt, err := s.db.Prepare("DELETE FROM daytask WHERE title = ? AND owner = ? AND date = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(taskName, taskOwner, taskDate)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil 
}

func (s *Storage) GetTaskForDay(taskOwner string, taskDate string) ([]storage.Task, error){
	const op = "storage.sqlite.GetTaskForDay"

	stmt, err := s.db.Prepare("SELECT * FROM daytask WHERE owner = ? AND date = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(taskOwner, taskDate)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var tasks []storage.Task

	for rows.Next() {
		var task storage.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Owner, &task.Date)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return tasks, nil 
}