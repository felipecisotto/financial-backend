package config

import (
	"sync"
)

type Event interface {
	EventName() string
}

type Handler interface {
	EventName() string
	Handle(Event)
}

type Publisher interface {
	RegisterHandler(h Handler)
	Publish(event Event)
}

type InMemoryPublisher struct {
	handlers map[string][]Handler
}

func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{
		handlers: make(map[string][]Handler),
	}
}

func (p *InMemoryPublisher) RegisterHandler(h Handler) {
	p.handlers[h.EventName()] = append(p.handlers[h.EventName()], h)
}

func (p *InMemoryPublisher) Publish(event Event) {
	for _, h := range p.handlers[event.EventName()] {
		go h.Handle(event)
	}
}

var (
	publisherInstance *InMemoryPublisher
	publisherOnce     sync.Once
)

// GetPublisher retorna a instância única do InMemoryPublisher
func GetPublisher() *InMemoryPublisher {
	publisherOnce.Do(func() {
		publisherInstance = NewInMemoryPublisher()
	})
	return publisherInstance
}
