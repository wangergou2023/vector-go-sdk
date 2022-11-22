package main

import (
	"context"
	"flag"
	sdk_wrapper "github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	sdk_wrapper.InitSDK(*serial)

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			stop <- true
			evtStreamHandler, _ := sdk_wrapper.Robot.Conn.EventStream(ctx,
				&vectorpb.EventRequest{})
			for {
				evt, _ := evtStreamHandler.Recv()
				print("FIRED!")
				evtRobotState := evt.Event.GetRobotState()
				if evtRobotState != nil {
					print(evtRobotState.String())
				}
			}
			return
		}
	}
}
