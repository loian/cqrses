package cqrses

import (
	"errors"
	"sync"
)

type EventBus interface {
	Publish (em EventMessage)
	Register(evtType EventType, handler EventHandler, lockExecution bool)
	Wait()
}

//commandHandlerCtrl is a control structure used to manage concurrency
//during the execution of a command handler
type evtMsgHandlerCtrl struct {
	evtHandler 	EventHandler
	lockExecution bool
	sync.Mutex
}

type eventBus struct {
	registry map[EventType][]*evtMsgHandlerCtrl
	sync.Mutex //mutex for the registry
	wg sync.WaitGroup
}

func (bus *eventBus) copyRegistry(et EventType) ([]*evtMsgHandlerCtrl, error) {
	if handlers, ok  := bus.registry[et]; ok {
		registryCopy := make([]*evtMsgHandlerCtrl, len(bus.registry[et]))
		registryCopy = append(registryCopy, handlers...)
		return registryCopy, nil
	} else {
		err := errors.New("xx")
		return nil, err
	}
}

func (bus *eventBus) Publish(em EventMessage) {
	defer bus.wg.Done()

	//mutex for accessing and copying the registry
	bus.Lock()
	handlers, err := bus.copyRegistry(em.Type())
	bus.Unlock()

	if err != nil { return }

	for _, h := range handlers {
		bus.wg.Add(1)
		//lock the callback
		if h.lockExecution {
			h.Lock()
		}

		go bus.executeHandler(h, em)
	}
}

func (bus *eventBus) executeHandler(eventController *evtMsgHandlerCtrl, em EventMessage) error{
	if eventController.lockExecution {
		defer eventController.Unlock()
	}

	return eventController.evtHandler.Handle(em)
}

func (bus *eventBus) Register(evtType EventType, handler EventHandler, lockExecution bool) {
	bus.Lock()
	defer bus.Unlock()
	bus.registry[evtType] = append(
		bus.registry[evtType],
		&evtMsgHandlerCtrl{
			handler,
			lockExecution,
			sync.Mutex{},
		})
}

func (bus *eventBus) Wait() {
	bus.wg.Wait()
}

func NewEventBus() *eventBus {
	return &eventBus{
		make(map[EventType][]*evtMsgHandlerCtrl),
		sync.Mutex{},
		sync.WaitGroup{},
	}
}
