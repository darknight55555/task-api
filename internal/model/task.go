package model

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskFilter struct {
	Done   *bool
	Limit  int
	Offset int
}
