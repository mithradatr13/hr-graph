package repository

import (
	"context"
	"database/sql"
	"fmt"
	"task-manager/internal/domain"
	"time"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `INSERT INTO tasks (title, status, assignee, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	return r.db.QueryRowContext(ctx, query, task.Title, task.Status, task.Assignee, task.CreatedAt, task.UpdatedAt).Scan(&task.ID)
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	query := `SELECT id, title, status, assignee, created_at, updated_at FROM tasks WHERE id = $1`
	var task domain.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Title, &task.Status, &task.Assignee, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	return &task, err
}

func (r *PostgresTaskRepository) List(ctx context.Context, status, assignee string, limit, offset int) ([]*domain.Task, error) {
	query := `SELECT id, title, status, assignee, created_at, updated_at FROM tasks WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if assignee != "" {
		query += fmt.Sprintf(" AND assignee = $%d", argIdx)
		args = append(args, assignee)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.Assignee, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `UPDATE tasks SET title = $1, status = $2, assignee = $3, updated_at = $4 WHERE id = $5`
	task.UpdatedAt = time.Now()
	res, err := r.db.ExecContext(ctx, query, task.Title, task.Status, task.Assignee, task.UpdatedAt, task.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
