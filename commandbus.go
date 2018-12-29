package cqrses

import (
	"errors"
	"sync"
)

type CommandBus interface {
	Dispatch (e CommandMessage)
	Register(evtType CommandType, cmdHandler CommandHandler, lockExecution bool)
}

//commandHandlerCtrl is a control structure used to manage concurrency
//during the execution of a command handler
type cmdMsgHandlerCtrl struct {
	cmdHandler 	CommandHandler
	lockExecution bool
	sync.Mutex
}

type commandBus struct {
	registry map[CommandType][]*cmdMsgHandlerCtrl
	sync.Mutex //mutex for the registry
	wg sync.WaitGroup
}

func (bus *commandBus) copyRegistry(ct CommandType) ([]*cmdMsgHandlerCtrl, error) {
	if handlers, ok  := bus.registry[ct]; ok {
		registryCopy := make([]*cmdMsgHandlerCtrl, len(bus.registry[ct]))
		registryCopy = append(registryCopy, handlers...)
		return registryCopy, nil
	} else {
		err := errors.New("xx")
		return nil, err
	}
}

func (bus *commandBus) Dispatch(cm CommandMessage) {
	defer bus.wg.Done()

	//mutex for accessing and copying the registry
	bus.Lock()
	handlers, err := bus.copyRegistry(cm.Type())
	bus.Unlock()

	if err != nil { return }

	for _, h := range handlers {
		bus.wg.Add(1)
		//lock the callback
		if h.lockExecution {
			h.Lock()
		}

		go bus.executeHandler(h, cm)
	}
}

func (bus *commandBus) executeHandler1(handlerController *cmdMsgHandlerCtrl, cm CommandMessage) {
	defer bus.wg.Done()

	if handlerController.lockExecution {
		defer handlerController.Unlock()
	}

	handlerController.cmdHandler.Handle(cm)
}

func (bus *commandBus) executeHandler(handlerController *cmdMsgHandlerCtrl, cm CommandMessage) error{
	if handlerController.lockExecution {
		defer handlerController.Unlock()
	}

	return handlerController.cmdHandler.Handle(cm)
}


func (bus *commandBus) Register(cmdType CommandType, handler CommandHandler, lockExecution bool) {
	bus.Lock()
	defer bus.Unlock()
	bus.registry[cmdType] = append(
		bus.registry[cmdType],
		&cmdMsgHandlerCtrl{
			handler,
			lockExecution,
			sync.Mutex{},
		})
}

func (bus *commandBus) Wait() {
	bus.wg.Wait()
}

func NewCommandBus() CommandBus {
	return &commandBus{
		make(map[CommandType][]*cmdMsgHandlerCtrl),
		sync.Mutex{},
		sync.WaitGroup{},
	}
}
