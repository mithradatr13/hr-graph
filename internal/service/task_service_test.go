package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task-manager/internal/domain"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskRepository) List(ctx context.Context, status, search string, limit, offset int) ([]*domain.Task, error) {
	args := m.Called(ctx, status, search, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}

type MockTaskCache struct {
	mock.Mock
}

func (m *MockTaskCache) Get(ctx context.Context, id int64) (*domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskCache) Set(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskCache) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTaskService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockCache := new(MockTaskCache)
	expectedTask := &domain.Task{ID: 1, Title: "Sample Task", Status: "pending"}

	mockCache.On("Get", mock.Anything, int64(1)).Return(nil, errors.New("cache miss"))
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedTask, nil)
	mockCache.On("Set", mock.Anything, expectedTask).Return(nil)

	taskService := NewTaskService(mockRepo, mockCache)
	task, err := taskService.GetTask(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, expectedTask.Title, task.Title)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestTaskService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockCache := new(MockTaskCache)

	mockCache.On("Get", mock.Anything, int64(99)).Return(nil, errors.New("cache miss"))
	mockRepo.On("GetByID", mock.Anything, int64(99)).Return(nil, errors.New("not found"))

	taskService := NewTaskService(mockRepo, mockCache)
	task, err := taskService.GetTask(context.Background(), 99)

	assert.Error(t, err)
	assert.Nil(t, task)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
