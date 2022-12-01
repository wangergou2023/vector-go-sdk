package main

import (
	"context"
	sdk_wrapper "github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
)

var Utterances = []string{"your name is", "you are"}
var Name = "Wirepod SDK Extesion Plugin"

func Action(transcribedText string, botSerial string) string {
	sdk_wrapper.InitSDK(botSerial)

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			sdk_wrapper.SayText("Ok. My name is Augustus")
			stop <- true
			return "intent_greeting_hello"
		}
	}
	return "intent_greeting_hello"
}
