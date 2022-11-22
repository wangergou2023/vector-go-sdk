package sdk_wrapper

import (
	"context"
	"fmt"
	"reflect"
)

type VectorEvent interface {
	Handle(ctx context.Context) error
}

type VectoEventRobotState interface {
}

type Dispatcher struct {
	events []string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Register(events ...VectorEvent) {
	for _, v := range events {
		d.events = append(d.events, reflect.TypeOf(v).String())
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, event VectorEvent) error {
	name := reflect.TypeOf(event).String()
	for _, v := range d.events {
		if v == name {
			return event.Handle(ctx)
		}
	}
	return fmt.Errorf("%s is not a registered event", name)
}

func (d *Dispatcher) StartListening() {
	/*
		ctx := context.Background()
		start := make(chan bool)
		stop := make(chan bool)

		for {
			select {
			case <-start:
				stop <- true
				evtStreamHandler, _ := Robot.Conn.EventStream(ctx,
					&vectorpb.EventRequest{})
				for {
					evt, _ := evtStreamHandler.Recv()
					evtRobotState := evt.Event.GetRobotState()
					if evtRobotState != nil {
						dis
					}
				}
				return
			}
		}
	*/
}
