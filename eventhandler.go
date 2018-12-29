package cqrses


type EventHandler interface {
	Handle(e EventMessage) error
}

