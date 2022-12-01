package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/images"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/weather"
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
			temperature := 16
			unit := pkg.GetTemperatureUnit()
			if unit == weather.WEATHER_UNIT_FARANHEIT {
				temperature = (temperature * 9 / 5) + 32
			}
			images.UseVectorEyeColorInImages(true)
			err1 := weather.DisplayTemperature(temperature, unit, 5000, true)
			if err1 != nil {
				println("ERROR " + err1.Error())
			}
			images.UseVectorEyeColorInImages(false)
			weather.DisplayCondition("shower rain", "09n", 5000, true)
			stop <- true
			return
		}
	}
}
