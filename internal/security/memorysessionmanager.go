package security

import "github.com/google/uuid"

type MemorySessionManager struct {
	sessions map[string]string
}

func NewMemorySessionManager() *MemorySessionManager {
	return &MemorySessionManager{sessions: make(map[string]string)}
}

func (m MemorySessionManager) CreateSession(user string) (string, error) {
	sessionID := uuid.New().String()
	m.sessions[user] = sessionID
	return sessionID, nil
}

func (m MemorySessionManager) GetSession(id string) (string, error) {
	for user, sessionID := range m.sessions {
		if sessionID == id {
			return user, nil
		}
	}
	return "", nil
}

func (m MemorySessionManager) TestSetSession(user string, sessionID string) {
	m.sessions[user] = sessionID
}
