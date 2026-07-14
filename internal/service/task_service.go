package service

import (
	"context"
	"strings"
	"task-api/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, title string) (model.Task, error)
	List(ctx context.Context, filter model.TaskFilter) ([]model.Task, error)
	GetByID(ctx context.Context, id int) (model.Task, error)
	Update(ctx context.Context, id int, title string, done bool) (model.Task, error)
	Delete(ctx context.Context, id int) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	taskService := TaskService{
		repo: repo,
	}

	return &taskService
}

func (t *TaskService) Create(ctx context.Context, title string) (model.Task, error) {
	if strings.TrimSpace(title) == "" {
		return model.Task{}, ErrInvalidTitle
	}

	task, err := t.repo.Create(ctx, title)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (t *TaskService) List(ctx context.Context, filter model.TaskFilter) ([]model.Task, error) {
	tasks, err := t.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *TaskService) GetByID(ctx context.Context, id int) (model.Task, error) {
	if id <= 0 {
		return model.Task{}, ErrInvalidID
	}

	task, err := t.repo.GetByID(ctx, id)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (t *TaskService) Update(ctx context.Context, id int, title string, done bool) (model.Task, error) {
	if strings.TrimSpace(title) == "" {
		return model.Task{}, ErrInvalidTitle
	} else if id <= 0 {
		return model.Task{}, ErrInvalidID
	}

	task, err := t.repo.Update(ctx, id, title, done)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (t *TaskService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	if err := t.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
