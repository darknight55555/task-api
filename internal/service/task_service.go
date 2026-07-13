package service

import (
	"strings"
	"task-api/internal/model"
)

type TaskRepository interface {
	Create(title string) model.Task
	List() []model.Task
	GetByID(id int) (model.Task, error)
	Update(id int, title string, done bool) (model.Task, error)
	Delete(id int) error
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

func (t *TaskService) Create(title string) (model.Task, error) {
	if strings.TrimSpace(title) == "" {
		return model.Task{}, ErrInvalidTitle
	}

	task := t.repo.Create(title)

	return task, nil
}

func (t *TaskService) List() []model.Task {
	tasks := t.repo.List()

	return tasks
}

func (t *TaskService) GetByID(id int) (model.Task, error) {
	if id <= 0 {
		return model.Task{}, ErrInvalidID
	}

	task, err := t.repo.GetByID(id)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (t *TaskService) Update(title string, id int, done bool) (model.Task, error) {
	if strings.TrimSpace(title) == "" {
		return model.Task{}, ErrInvalidTitle
	} else if id <= 0 {
		return model.Task{}, ErrInvalidID
	}

	task, err := t.repo.Update(id, title, done)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (t *TaskService) Delete(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	if err := t.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
