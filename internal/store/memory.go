package store

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/pablo-banker/taskforge/internal/task"
)

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyExists  = errors.New("task already exists")
	ErrTaskNotCancellable = errors.New("task cannot be cancelled")
)

type MemoryStore struct {
	mu       sync.RWMutex
	tasks    map[string]task.Task
	watchers map[string]map[chan task.Task]struct{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:    make(map[string]task.Task),
		watchers: make(map[string]map[chan task.Task]struct{}),
	}
}

func (s *MemoryStore) Create(_ context.Context, t task.Task) (task.Task, error) {
	s.mu.Lock()

	if _, exists := s.tasks[t.ID]; exists {
		s.mu.Unlock()
		return task.Task{}, ErrTaskAlreadyExists
	}

	s.tasks[t.ID] = t
	s.mu.Unlock()

	s.notify(t)

	return t, nil
}

func (s *MemoryStore) Get(_ context.Context, id string) (task.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, exists := s.tasks[id]
	if !exists {
		return task.Task{}, ErrTaskNotFound
	}

	return t, nil
}

func (s *MemoryStore) List(_ context.Context, queueName string, status task.Status) ([]task.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]task.Task, 0, len(s.tasks))

	for _, t := range s.tasks {
		if queueName != "" && t.Queue != queueName {
			continue
		}

		if status != "" && t.Status != status {
			continue
		}

		result = append(result, t)
	}

	return result, nil
}

func (s *MemoryStore) ReadyScheduled(now time.Time) []task.Task {
	s.mu.Lock()

	ready := make([]task.Task, 0)

	for id, t := range s.tasks {
		if t.Status != task.StatusScheduled {
			continue
		}

		if t.RunAt.After(now) {
			continue
		}

		t.Status = task.StatusPending
		t.UpdatedAt = now
		t.LastError = ""

		s.tasks[id] = t
		ready = append(ready, t)
	}

	s.mu.Unlock()

	for _, t := range ready {
		s.notify(t)
	}

	return ready
}

func (s *MemoryStore) StartTask(id string) (task.Task, bool, error) {
	now := time.Now()

	s.mu.Lock()

	t, exists := s.tasks[id]
	if !exists {
		s.mu.Unlock()
		return task.Task{}, false, ErrTaskNotFound
	}

	if t.Status != task.StatusPending {
		s.mu.Unlock()
		return t, false, nil
	}

	t.Status = task.StatusRunning
	t.Attempts++
	t.UpdatedAt = now

	s.tasks[id] = t
	s.mu.Unlock()

	s.notify(t)

	return t, true, nil
}

func (s *MemoryStore) CompleteTask(id string) (task.Task, error) {
	now := time.Now()

	s.mu.Lock()

	t, exists := s.tasks[id]
	if !exists {
		s.mu.Unlock()
		return task.Task{}, ErrTaskNotFound
	}

	t.Status = task.StatusCompleted
	t.UpdatedAt = now
	t.LastError = ""

	s.tasks[id] = t
	s.mu.Unlock()

	s.notify(t)

	return t, nil
}

func (s *MemoryStore) FailTask(id string, message string, nextRunAt time.Time) (task.Task, error) {
	now := time.Now()

	s.mu.Lock()

	t, exists := s.tasks[id]
	if !exists {
		s.mu.Unlock()
		return task.Task{}, ErrTaskNotFound
	}

	t.LastError = message
	t.UpdatedAt = now

	if t.Attempts >= t.MaxAttempts {
		t.Status = task.StatusDeadLetter
	} else {
		t.Status = task.StatusScheduled
		t.RunAt = nextRunAt
	}

	s.tasks[id] = t
	s.mu.Unlock()

	s.notify(t)

	return t, nil
}

func (s *MemoryStore) CancelTask(id string) (task.Task, error) {
	now := time.Now()

	s.mu.Lock()

	t, exists := s.tasks[id]
	if !exists {
		s.mu.Unlock()
		return task.Task{}, ErrTaskNotFound
	}

	if t.Status == task.StatusRunning || t.IsTerminal() {
		s.mu.Unlock()
		return task.Task{}, ErrTaskNotCancellable
	}

	t.Status = task.StatusCancelled
	t.UpdatedAt = now

	s.tasks[id] = t
	s.mu.Unlock()

	s.notify(t)

	return t, nil
}

func (s *MemoryStore) Subscribe(id string) (<-chan task.Task, func()) {
	ch := make(chan task.Task, 16)

	s.mu.Lock()

	if _, exists := s.watchers[id]; !exists {
		s.watchers[id] = make(map[chan task.Task]struct{})
	}

	s.watchers[id][ch] = struct{}{}

	s.mu.Unlock()

	cancel := func() {
		s.mu.Lock()

		if watchers, exists := s.watchers[id]; exists {
			delete(watchers, ch)

			if len(watchers) == 0 {
				delete(s.watchers, id)
			}
		}

		close(ch)

		s.mu.Unlock()
	}

	return ch, cancel
}

func (s *MemoryStore) notify(t task.Task) {
	s.mu.RLock()

	watchers := make([]chan task.Task, 0)

	for ch := range s.watchers[t.ID] {
		watchers = append(watchers, ch)
	}

	s.mu.RUnlock()

	for _, ch := range watchers {
		select {
		case ch <- t:
		default:
		}
	}
}
