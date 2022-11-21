package main

import (
	"context"
	"flag"
	"fmt"
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
			sdk_wrapper.FaceEnrollmentStart("Ambrogio", 99)
			faces := sdk_wrapper.FaceEnrollmentListAll()
			println("")
			println("KNOWN FACES")
			for i := 0; i < len(faces); i++ {
				println(fmt.Sprintf("%d) %s", faces[i].FaceId, faces[i].Name))
			}
			stop <- true
			return
		}
	}
}
