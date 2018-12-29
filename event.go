package cqrses

import "reflect"

type EventType string

type EventMessage interface {
	AggregateID() AggregateID
	Event() interface{}
	Type() EventType
}

type EventMessageBase struct {
	aggregateID AggregateID
	event interface{}
}

func (emb *EventMessageBase) AggregateId() AggregateID {
	return emb.aggregateID
}

func (emb *EventMessageBase) Event() interface{} {
	return emb.event
}

func (emb *EventMessageBase) Type() EventType {
	return EventType(reflect.TypeOf(emb.event).String())
}