package repository

import (
	"context"
	"sync"
	"task-api/internal/model"
	"time"
)

type MemoryTaskRepository struct {
	tasks  map[int]model.Task
	nextID int
	mu     sync.RWMutex
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	taskRepository := MemoryTaskRepository{
		tasks:  make(map[int]model.Task),
		nextID: 1,
	}

	return &taskRepository
}

func (m *MemoryTaskRepository) Create(ctx context.Context, title string) (model.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task := model.Task{
		ID:        m.nextID,
		CreatedAt: time.Now(),
		Title:     title,
	}

	m.nextID += 1

	m.tasks[task.ID] = task

	return m.tasks[task.ID], nil
}

func (m *MemoryTaskRepository) List(ctx context.Context) ([]model.Task, error) {
	list := []model.Task{}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, v := range m.tasks {
		list = append(list, v)
	}

	return list, nil
}

func (m *MemoryTaskRepository) GetByID(ctx context.Context, id int) (model.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, ok := m.tasks[id]
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (m *MemoryTaskRepository) Update(ctx context.Context, id int, title string, done bool) (model.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[id]
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}
	task.Title = title
	task.Done = done

	m.tasks[id] = task

	return task, nil
}

func (m *MemoryTaskRepository) Delete(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.tasks[id]
	if !ok {
		return ErrTaskNotFound
	}

	delete(m.tasks, id)

	return nil
}
