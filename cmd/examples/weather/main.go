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
			err1 := sdk_wrapper.DisplayTemperature(16, sdk_wrapper.WEATHER_UNIT_CELSIUS, 5000, true)
			if err1 != nil {
				println("ERROR " + err1.Error())
			}
			sdk_wrapper.DisplayCondition("shower rain", "09n", 5000, true)
			stop <- true
			return
		}
	}
}
