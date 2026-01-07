package model

import "time"

type StatusTask string

type Task struct {
	ID          uint64     `json:"id"`
	Description string     `json:"description"`
	Status      StatusTask `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
