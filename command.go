package cqrses

import "reflect"

type CommandType string

type CommandMessage interface {
	AggregateID() AggregateID
	Command() interface{}
	Type() CommandType
}


type CommandMessageBase struct {
	aggregateID AggregateID
	command interface{}
}

func (cm CommandMessageBase) AggregateID() AggregateID {
	return cm.aggregateID
}

func (cm CommandMessageBase) Command() interface{} {
	return cm.command
}

func (cm CommandMessageBase) Type() CommandType {
	return CommandType(reflect.TypeOf(cm.command).String())
}


func NewCommandMessage(id string, command interface{}) CommandMessage {
	return CommandMessageBase{
		aggregateID: AggregateID(id),
		command: command,
	}
}