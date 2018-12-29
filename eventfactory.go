package cqrses

import (
	"fmt"
)

type EventFactory interface {
	Get(et EventType) (interface{}, error)
}

type eventFactory struct {
	delegates map[EventType]func() interface{}
}

func NewEventFactory() *eventFactory{
	return &eventFactory{
		delegates: make(map[EventType]func() interface{}),
	}
}

func (ef *eventFactory) Register(et EventType, delegate func() interface{}) error{
	if _, ok := ef.delegates[et]; ok {
		return fmt.Errorf("delegate already register for event type %s", et)
	}

	ef.delegates[et] = delegate
	return nil
}

func (ef *eventFactory) Get(et EventType) (interface{}, error) {
	if _, ok := ef.delegates[et]; !ok {
		return nil, fmt.Errorf("no delegate registered for event type %s", et)
	}

	return ef.delegates[et](), nil
}