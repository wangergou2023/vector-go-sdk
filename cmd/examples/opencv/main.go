package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	sdk_wrapper "github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"log"
	"os/exec"
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
			doOpenCVStuff(50)
			stop <- true
			return
		}
	}
}

func doOpenCVStuff(numSteps int) {
	for i := 0; i <= numSteps; i++ {
		fName := fmt.Sprintf("/tmp/camera%02d.jpg", i)
		err := sdk_wrapper.SaveHiResCameraPicture(fName)
		if err == nil {
			if err == nil {
				cmd := exec.Command("python", "hand_detection.py", fName)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("FRAME %d, Output: %s\n", i, out.String())
				sdk_wrapper.SayText(out.String())
				time.Sleep(time.Duration(2000) * time.Millisecond)
			} else {
				println("OPENCV Python script not found!")
			}
		}
	}

}
