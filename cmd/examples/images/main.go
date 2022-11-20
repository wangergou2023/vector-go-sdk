package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
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
			sdk_wrapper.WriteText("HI", 32, true, 5000, false)
			sdk_wrapper.SayText("Hi")
			sdk_wrapper.WriteText("I am Vector", 32, false, 5000, false)
			sdk_wrapper.SayText("I am Vector")
			sdk_wrapper.DisplayImage("data/images/birthday-cake.png", 5000, false)
			sdk_wrapper.SayText("Happy birthday!")
			stop <- true
			return
		}
	}
}
