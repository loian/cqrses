package cqrses

import (
	"fmt"
	"reflect"
	"sync"
)

type AggregateID string

type AggregateRoot interface {
	AggregateID() AggregateID
	Version() int64
	InitialVersion() int64
	IncrementVersion()
	Changes() []EventMessage
	AppendChange(eventMessage EventMessage)
	ClearChanges()
	Apply(eventMessage EventMessage) error

}

type EventProcessorAlreadyDefined struct{
	message string
}

func (e EventProcessorAlreadyDefined) Error() string{
	return e.message
}

type AggregateRootBase struct {
	id AggregateID
	version int64
	changes []EventMessage
	changesMutex sync.Mutex
	initialVersion int64

	registry map[EventType] func(message EventMessage)
}

func (arb *AggregateRootBase) AggregateID() AggregateID {
	return arb.id
}

func (arb *AggregateRootBase) Version() int64 {
	return arb.version //+ len(events)
}

func (arb *AggregateRootBase) InitialVersion() int64 {
	return arb.initialVersion
}

func (arb *AggregateRootBase) IncrementVersion() {
	arb.version++
}

func (arb *AggregateRootBase) Changes() []EventMessage {
	return arb.changes
}

func (arb *AggregateRootBase) AppendChange(event EventMessage) {
	arb.changesMutex.Lock()
	defer arb.changesMutex.Unlock()

	arb.changes = append(arb.changes, event)
	arb.IncrementVersion()
}

func (arb *AggregateRootBase) ClearChanges() {
	arb.changes = []EventMessage{}
}

func (arb *AggregateRootBase) AddApplyProcessor (et EventType, f func (message EventMessage)) error {
	if _, ok := arb.registry[et]; ok {
		return EventProcessorAlreadyDefined{
			fmt.Sprintf(`apply processor for "%s" event type already defined`, string(et)),
		}
	}

	arb.registry[et] = f
	return nil
}

func (arb *AggregateRootBase) Apply(eventMessage EventMessage) error{
	eventType := EventType(reflect.TypeOf(eventMessage).String())

	proc, ok := arb.registry[eventType]
	if !ok {
		return EventProcessorAlreadyDefined{
			fmt.Sprintf(`no apply processor for event type "%s"`, string(eventType)),
		}
	}

	proc(eventMessage)
	return nil
}

