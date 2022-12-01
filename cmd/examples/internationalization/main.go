package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/voice"
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
			pkg.SetRobotName("Augustus")
			voice.SetLanguage(voice.LANGUAGE_ENGLISH)
			voice.SayText("Hello world!")
			voice.SetLanguage(voice.LANGUAGE_ITALIAN)
			voice.SayText("Ciao mondo!")
			voice.SetLanguage(voice.LANGUAGE_SPANISH)
			voice.SayText(voice.Translate("Hello, world!", voice.LANGUAGE_ENGLISH, voice.LANGUAGE_SPANISH))
			voice.SetLanguage(voice.LANGUAGE_FRENCH)
			voice.SayText(voice.Translate("Hello, world!", voice.LANGUAGE_ENGLISH, voice.LANGUAGE_FRENCH))
			voice.SetLanguage(voice.LANGUAGE_GERMAN)
			voice.SayText(voice.Translate("Hello, world!", voice.LANGUAGE_ENGLISH, voice.LANGUAGE_GERMAN))
			voice.SetLanguage(voice.LANGUAGE_JAPANESE)
			voice.SayText(voice.Translate("Hello, world!", voice.LANGUAGE_ENGLISH, voice.LANGUAGE_JAPANESE))
			stop <- true
			return
		}
	}
}
