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

	stmt, err := db.Prepare( //TODO: добавить енум не сделано в процессе завершено 
		`CREATE TABLE IF NOT EXISTS daytask(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		owner TEXT NOT NULL,
		date DATE NOT NULL,
		status TEXT, 
		type TEXT);
		CREATE INDEX IF NOT EXIST idx_id ON daytask(id);
		`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	
	}
	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(  
		`CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL);
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

func (s *Storage) SaveTask(taskName string, taskDescription string, taskOwner string, taskDate string, taskStatus string, taskType string) (int64, error){
	const op = "storage.sqlite.SaveTask"

	stmt, err := s.db.Prepare("INSERT INTO daytask(title, description, owner, date, status, type) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(taskName, taskDescription, taskOwner, taskDate, taskStatus, taskType)
	if err != nil { 
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil 
}

func (s *Storage) DeleteTask(id int64) (error){
	const op = "storage.sqlite.DeleteTask"

	stmt, err := s.db.Prepare("DELETE FROM daytask WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
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
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Owner, &task.Date, &task.Status, &task.Type)

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

func (s *Storage) GetAllTasks(taskOwner string) ([]storage.Task, error){
	const op = "storage.sqlite.GetTaskForDay"

	stmt, err := s.db.Prepare("SELECT * FROM daytask WHERE owner = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(taskOwner)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var tasks []storage.Task

	for rows.Next() {
		var task storage.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Owner, &task.Date, &task.Status, &task.Type)

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

func (s *Storage) UpdateTask(taskID int64, taskName string, taskDescription string, taskOwner string, taskDate string, taskStatus string, taskType string) (error){
	const op = "storage.sqlite.UpdateTask"

	stmt, err := s.db.Prepare("UPDATE daytask SET title = ?, description = ?, owner = ?, date = ?, status = ?, type = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(taskName, taskDescription, taskOwner, taskDate, taskStatus, taskType, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil 
}

func (s *Storage) CreateUser(username string, password string) (int64, error){
	const op = "storage.sqlite.CreateUser"

	stmt, err := s.db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(username, password)
	if err != nil { 
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil 
}