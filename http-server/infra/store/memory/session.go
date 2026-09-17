package memory

import (
	"neomatica/neosync/infra/logger"
	"sync"

	"github.com/google/uuid"
)

type SessionService struct {
	sync.Mutex
	sessions map[string]chan []byte
}

func NewSessionService() *SessionService {
	return &SessionService{sessions: make(map[string]chan []byte)}
}

func (s *SessionService) Register(requestId string, ch chan []byte) {
	s.Lock()
	defer s.Unlock()

	s.sessions[requestId] = ch
}

func (s *SessionService) Unregister(requestId string) {
	s.Lock()
	defer s.Unlock()

	delete(s.sessions, requestId)
}

func (s *SessionService) Dispatch(buffer []byte) {
	if len(buffer) < 19 {
		logger.Error("Sessino Dispatch: Failed to dispatch buffer")
		return
	}

	reqIdBytes := buffer[3:19]
	requestID := uuid.UUID{}
	copy(requestID[:], reqIdBytes)

	s.Lock()
	ch, exists := s.sessions[requestID.String()]
	s.Unlock()

	if exists {
		select {
		case ch <- buffer:
			/* OK */
		default:
			logger.Error("Session Dispatch: Channel blocked for requestId={%s}", requestID.String())
		}
	}
}
