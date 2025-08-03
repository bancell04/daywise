package models

import "time"

type Task struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}
