package main

import (
	"context"
	"flag"
	sdk_wrapper "github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/audio"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/voice"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	var sentence = flag.String("sentence", "", "Sentence to say")
	var audioFile = flag.String("audiofile", "", "Audio file to stream")
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
			println("OK")
			if *sentence != "" {
				voice.SayText(*sentence)
			}
			if *audioFile != "" {
				ret := audio.PlaySound(*audioFile)
				println(ret)
			}
			stop <- true
			return
		}
	}
}
