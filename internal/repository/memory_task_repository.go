package repository

import (
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

func (m *MemoryTaskRepository) Create(title string) model.Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	task := model.Task{
		ID:        m.nextID,
		CreatedAt: time.Now(),
		Title:     title,
	}

	m.nextID += 1

	m.tasks[task.ID] = task

	return m.tasks[task.ID]
}

func (m *MemoryTaskRepository) List() []model.Task {
	list := []model.Task{}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, v := range m.tasks {
		list = append(list, v)
	}

	return list
}

func (m *MemoryTaskRepository) GetByID(id int) (model.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, ok := m.tasks[id]
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (m *MemoryTaskRepository) Update(id int, title string, done bool) (model.Task, error) {
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

func (m *MemoryTaskRepository) Delete(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.tasks[id]
	if !ok {
		return ErrTaskNotFound
	}

	delete(m.tasks, id)

	return nil
}
