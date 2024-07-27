package storage

import "errors"

var (
	ErrLoginNotFound = errors.New("login not found")
	ErrLoginExists   = errors.New("login exists")
	ErrWrongPassword   = errors.New("wrong password")
	ErrIncorrectDate   = errors.New("incorrect date")
)

type Task struct {
	ID          int64    `json:"id"`
	Title       string 	 `json:"title"`
	Description string 	 `json:"description"`
	Owner		string 	 `json:"owner"`
	Date        string 	 `json:"date"`
	Status 		string 	 `json:"status,omitempty"`	
	Type 		string 	 `json:"type,omitempty"`
}