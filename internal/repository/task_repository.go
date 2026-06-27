// Package repository provides PostgreSQL implementations of domain repositories.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/juanozorio/task-api/internal/domain"
)

type taskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new PostgreSQL-backed TaskRepository.
func NewTaskRepository(db *sql.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, description, is_completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		task.ID, task.Description, task.IsCompleted, task.CreatedAt, task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	task.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE tasks
		SET description = $1, is_completed = $2, updated_at = $3
		WHERE id = $4
	`
	result, err := r.db.ExecContext(ctx, query,
		task.Description, task.IsCompleted, task.UpdatedAt, task.ID,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update task rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

func (r *taskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `
		SELECT id, description, is_completed, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	task := &domain.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return task, nil
}

func (r *taskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	query := `
		SELECT id, description, is_completed, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`
	return r.scanTasks(ctx, query)
}

func (r *taskRepository) GetByStatus(ctx context.Context, isCompleted bool) ([]*domain.Task, error) {
	query := `
		SELECT id, description, is_completed, created_at, updated_at
		FROM tasks
		WHERE is_completed = $1
		ORDER BY created_at DESC
	`
	return r.scanTasks(ctx, query, isCompleted)
}

func (r *taskRepository) scanTasks(ctx context.Context, query string, args ...any) ([]*domain.Task, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		if err := rows.Scan(
			&task.ID, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if tasks == nil {
		tasks = []*domain.Task{}
	}

	return tasks, nil
}
