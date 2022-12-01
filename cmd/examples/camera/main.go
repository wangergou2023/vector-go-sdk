package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/animations"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/camera"
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
			animations.PlayAnimation("anim_generic_look_up_01", 0, false, false, false)
			camera.SaveHiResCameraPicture("camera.jpg")
			stop <- true
			return
		}
	}
}
