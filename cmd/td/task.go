package main

import "time"

type Task struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	ID        int       `json:"id"`
	IsDone    bool      `json:"is_done"`
}
