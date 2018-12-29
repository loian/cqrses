package cqrses


type CommandHandler interface {
	Handle(e CommandMessage) error
}

