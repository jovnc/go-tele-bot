package state

import (
	"sync"
)

type UserState struct {
	Step     string // "select_options", "enter_email", "done"
	Selected map[string]bool
	Email    string
}

// Manager manages conversation states for all users
type Manager struct {
	mu     sync.RWMutex
	states map[int64]*UserState
}

// NewManager creates a new Manager instance
func NewManager() *Manager {
	return &Manager{
		states: make(map[int64]*UserState),
	}
}

func (m *Manager) Get(userID int64) *UserState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.states[userID]
}

func (m *Manager) Set(userID int64, s *UserState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[userID] = s
}

func (m *Manager) Delete(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, userID)
}

