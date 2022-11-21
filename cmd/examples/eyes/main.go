package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"time"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	var hue = flag.String("hue", "", "Hue")
	var saturation = flag.String("saturation", "", "Saturation")
	flag.Parse()

	sdk_wrapper.InitSDK(*serial)

	if *hue == "" {
		*hue = "2"
	}
	if *saturation == "" {
		*hue = "3"
	}

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			sdk_wrapper.SetCustomEyeColor(*hue, *saturation)
			time.Sleep(time.Duration(5000) * time.Millisecond)
			stop <- true
			return
		}
	}
}
