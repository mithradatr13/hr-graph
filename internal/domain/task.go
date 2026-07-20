package domain

import (
	"context"
	"time"
)

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Assignee  string    `json:"assignee,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id int64) (*Task, error)
	List(ctx context.Context, status, assignee string, limit, offset int) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int64) error
}

type TaskCache interface {
	Get(ctx context.Context, id int64) (*Task, error)
	Set(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int64) error
}
