package service

import (
	"context"
	"task-manager/internal/domain"

	"github.com/prometheus/client_golang/prometheus"
)

var TasksCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "tasks_count",
	Help: "Current number of tasks in the system",
})

type TaskService struct {
	repo  domain.TaskRepository
	cache domain.TaskCache
}

func NewTaskService(repo domain.TaskRepository, cache domain.TaskCache) *TaskService {
	return &TaskService{repo: repo, cache: cache}
}

func (s *TaskService) CreateTask(ctx context.Context, task *domain.Task) error {
	err := s.repo.Create(ctx, task)
	if err == nil {
		TasksCountGauge.Inc()
		if s.cache != nil {
			_ = s.cache.Set(ctx, task)
		}
	}
	return err
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*domain.Task, error) {
	if s.cache != nil {
		task, err := s.cache.Get(ctx, id)
		if err == nil && task != nil {
			return task, nil
		}
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, task)
	}
	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, status, assignee string, limit, offset int) ([]*domain.Task, error) {
	return s.repo.List(ctx, status, assignee, limit, offset)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domain.Task) error {
	err := s.repo.Update(ctx, task)
	if err == nil && s.cache != nil {
		_ = s.cache.Delete(ctx, task.ID) // Cache Invalidation
	}
	return err
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		TasksCountGauge.Dec()
		if s.cache != nil {
			_ = s.cache.Delete(ctx, id) // Cache Invalidation
		}
	}
	return err
}
