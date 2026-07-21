package service

import (
	"context"
	"log/slog"
	"task-manager/internal/domain"

	"github.com/prometheus/client_golang/prometheus"
)

var TasksCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "tasks_count",
	Help: "Current number of tasks in the system",
})

type TaskService struct {
	repo   domain.TaskRepository
	cache  domain.TaskCache
	logger *slog.Logger
}

func NewTaskService(repo domain.TaskRepository, cache domain.TaskCache, logger *slog.Logger) *TaskService {
	if logger == nil {
		logger = slog.Default()
	}
	return &TaskService{repo: repo, cache: cache, logger: logger}
}

func (s *TaskService) CreateTask(ctx context.Context, task *domain.Task) error {
	err := s.repo.Create(ctx, task)
	if err != nil {
		return err
	}

	TasksCountGauge.Inc()

	if s.cache != nil {
		if cacheErr := s.cache.Set(ctx, task); cacheErr != nil {
			s.logger.Warn("failed to set cache for task", "task_id", task.ID, "error", cacheErr)
		}
	}
	return nil
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*domain.Task, error) {
	if s.cache != nil {
		task, err := s.cache.Get(ctx, id)
		if err == nil && task != nil {
			return task, nil
		}
		if err != nil {
			s.logger.Warn("redis get error for task", "task_id", id, "error", err)
		}
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if cacheErr := s.cache.Set(ctx, task); cacheErr != nil {
			s.logger.Warn("failed to update cache for task", "task_id", task.ID, "error", cacheErr)
		}
	}

	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, status, assignee string, limit, offset int) ([]*domain.Task, error) {
	return s.repo.List(ctx, status, assignee, limit, offset)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domain.Task) error {
	err := s.repo.Update(ctx, task)
	if err != nil {
		return err
	}

	if s.cache != nil {
		if cacheErr := s.cache.Delete(ctx, task.ID); cacheErr != nil {
			s.logger.Warn("failed to invalidate cache for task", "task_id", task.ID, "error", cacheErr)
		}
	}
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	TasksCountGauge.Dec()

	if s.cache != nil {
		if cacheErr := s.cache.Delete(ctx, id); cacheErr != nil {
			s.logger.Warn("failed to invalidate cache for task", "task_id", id, "error", cacheErr)
		}
	}
	return nil
}
