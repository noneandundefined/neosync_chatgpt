package event

import (
	"neomatica/neosync-tcp/pkg"
	"sync"
)

type EventEmitter struct {
	listeners     map[string][]func(any)
	onceListeners map[string][]func(any)
	mu            sync.RWMutex
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners:     make(map[string][]func(any)),
		onceListeners: make(map[string][]func(any)),
	}
}

func (e *EventEmitter) On(event string, handler func(any)) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.listeners[event] = append(e.listeners[event], handler)
}

func (e *EventEmitter) Once(event string, handler func(any)) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.onceListeners[event] = append(e.onceListeners[event], handler)
}

func (e *EventEmitter) Emit(event string, data any) {
	e.mu.RLock()
	handlers := append([]func(any){}, e.listeners[event]...)
	onceHandlers := append([]func(any){}, e.onceListeners[event]...)
	e.mu.RUnlock()

	if len(onceHandlers) > 0 {
		e.mu.Lock()
		delete(e.onceListeners, event)
		e.mu.Unlock()
	}

	for _, handler := range handlers {
		h := handler
		pkg.RecoverGo("event-"+event, func() {
			h(data)
		})
	}

	for _, handler := range onceHandlers {
		h := handler
		pkg.RecoverGo("event-"+event+"-once", func() {
			h(data)
		})
	}
}
