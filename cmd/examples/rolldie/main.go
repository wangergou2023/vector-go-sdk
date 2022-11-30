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
			sdk_wrapper.MoveHead(3.0)
			sdk_wrapper.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, sdk_wrapper.IMAGE_TRANSITION_FADE_IN, 10)
			sdk_wrapper.DisplayAnimatedGif("data/images/dice/roll-the-dice.gif", sdk_wrapper.ANIMATED_GIF_SPEED_NORMAL, 3, false)
			stop <- true
			return
		}
	}
}
