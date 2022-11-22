package main

import (
	"context"
	"flag"
	sdk_wrapper "github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"time"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	sdk_wrapper.InitSDK(*serial)

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)
	event := make(chan *vectorpb.Event)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop, event)
	}()

	for {
		select {
		case <-start:
			//evtStreamHandler, _ := sdk_wrapper.Robot.Conn.EventStream(ctx, &vectorpb.EventRequest{})
			for {
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
			stop <- true
			return
		case <-event:
			evt := <-event
			evtRobotState := evt.GetRobotState()
			if evtRobotState != nil {
				if evtRobotState.TouchData.IsBeingTouched {
					println("You touch me!")
				}
			}
		}
	}
}
