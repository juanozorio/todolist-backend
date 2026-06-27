// Package domain contains the core business models and repository interfaces.
package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Sentinel errors for domain-level error handling.
var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrInvalidPayload = errors.New("invalid payload")
)

// Task represents a task entity.
type Task struct {
	ID          string    `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	IsCompleted bool      `json:"isCompleted" db:"is_completed"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateTaskInput holds the data required to create a new task.
type CreateTaskInput struct {
	Description string `json:"description" validate:"required,min=1,max=500"`
}

// UpdateTaskInput holds the data required to update an existing task.
type UpdateTaskInput struct {
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=500"`
	IsCompleted *bool   `json:"isCompleted,omitempty"`
}

// NewTask creates a new Task with a generated UUID.
func NewTask(input CreateTaskInput) *Task {
	now := time.Now().UTC()
	return &Task{
		ID:          uuid.New().String(),
		Description: input.Description,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// TaskRepository defines the persistence contract for tasks.
type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	GetAll(ctx context.Context) ([]*Task, error)
	GetByStatus(ctx context.Context, isCompleted bool) ([]*Task, error)
}
