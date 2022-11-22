package main

import (
	"context"
	"flag"
	sdk_wrapper "github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"time"
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
			go func() {
				println("Listening for events...")
				for {
					evt := sdk_wrapper.WaitForEvent()
					if evt != nil {
						evtRobotState := evt.GetRobotState()
						if evtRobotState != nil {
							if evtRobotState.TouchData.IsBeingTouched == true {
								println("I am being touched.")
							}
						}
					}
				}
			}()
			for {
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
			stop <- true
			return
		}
	}
}
