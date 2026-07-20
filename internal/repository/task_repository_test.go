package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"task-manager/internal/domain"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	return db, mock
}

func TestPostgresTaskRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresTaskRepository(db)

	now := time.Now()
	task := &domain.Task{
		Title:     "Test Task",
		Status:    "pending",
		Assignee:  "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tasks (title, status, assignee, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`)).
		WithArgs("Test Task", "pending", "", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.Create(context.Background(), task)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresTaskRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresTaskRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "title", "status", "assignee", "created_at", "updated_at"}).
		AddRow(1, "Test Task", "pending", "", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, status, assignee, created_at, updated_at FROM tasks WHERE id = $1`)).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	task, err := repo.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "Test Task", task.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}
