package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
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
			println("ANIMATION LIST:")
			aList := sdk_wrapper.LoadAnimationList()
			for i := 0; i < len(aList); i++ {
				println(aList[i])
			}
			sdk_wrapper.PlayAnimation("anim_weather_sunny_01", 1, false, false, false)
			time.Sleep(time.Duration(5000) * time.Millisecond)
			stop <- true
			return
		}
	}
}
