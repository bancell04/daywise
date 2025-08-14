package models

import "time"

type Task struct {
	ID       *int      `json:"id"`
	Title    string    `json:"title"`
	Category int       `json:"category"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}
