package antispam

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"sync"
	"time"
)

type state struct {
	ErrorsCount  int
	LastError    time.Time
	BlockedUntil time.Time
}

type Manager struct {
	mu    sync.Mutex
	users map[string]*state
}

func New() *Manager {
	return &Manager{
		users: make(map[string]*state),
	}
}

func (m *Manager) IsBlocked(uuid string) (bool, time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	st, ok := m.users[uuid]
	if !ok {
		return false, 0
	}

	if time.Now().Before(st.BlockedUntil) {
		return true, time.Until(st.BlockedUntil)
	}

	return false, 0
}

func (m *Manager) RegisterError(uuid string) time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	st, ok := m.users[uuid]
	if !ok {
		st = &state{}
		m.users[uuid] = st
	}

	if now.Sub(st.LastError) > time.Minute {
		st.ErrorsCount = 0
	}

	st.LastError = now
	st.ErrorsCount++

	if st.ErrorsCount >= 3 {
		idx := st.ErrorsCount - 3
		if idx >= len(constants.PenaltyDurations) {
			idx = len(constants.PenaltyDurations) - 1
		}

		block := constants.PenaltyDurations[idx]
		st.BlockedUntil = now.Add(block)

		logger.Info("Antispam: user={%s} errors={%d} block={%v} until={%v}", uuid, st.ErrorsCount, block, st.BlockedUntil)

		return block
	}

	logger.Info("Antispam: user={%s} errors={%d} (no block yet)", uuid, st.ErrorsCount)

	return 0
}

func (m *Manager) Reset(uuid string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.users, uuid)
}
