// Package service contains business logic for the task domain.
package service

import (
	"context"
	"fmt"

	"github.com/juanozorio/task-api/internal/domain"
)

// TaskService defines the use-case operations for tasks.
type TaskService interface {
	Create(ctx context.Context, input domain.CreateTaskInput) (*domain.Task, error)
	Update(ctx context.Context, id string, input domain.UpdateTaskInput) (*domain.Task, error)
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	GetAll(ctx context.Context) ([]*domain.Task, error)
	GetByStatus(ctx context.Context, isCompleted bool) ([]*domain.Task, error)
}

type taskService struct {
	repo domain.TaskRepository
}

// NewTaskService creates a TaskService with the given repository.
func NewTaskService(repo domain.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) Create(ctx context.Context, input domain.CreateTaskInput) (*domain.Task, error) {
	if input.Description == "" {
		return nil, fmt.Errorf("%w: description is required", domain.ErrInvalidPayload)
	}

	task := domain.NewTask(input)
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

func (s *taskService) Update(ctx context.Context, id string, input domain.UpdateTaskInput) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Description != nil {
		if *input.Description == "" {
			return nil, fmt.Errorf("%w: description cannot be empty", domain.ErrInvalidPayload)
		}
		task.Description = *input.Description
	}
	if input.IsCompleted != nil {
		task.IsCompleted = *input.IsCompleted
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	return task, nil
}

func (s *taskService) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *taskService) GetAll(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *taskService) GetByStatus(ctx context.Context, isCompleted bool) ([]*domain.Task, error) {
	return s.repo.GetByStatus(ctx, isCompleted)
}
