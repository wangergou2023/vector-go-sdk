package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/audio"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/images"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/motors"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/voice"
	"image/color"
	"math/rand"
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
			motors.MoveHead(3.0)
			images.SetBackgroundColor(color.RGBA{0, 0, 0, 0})

			images.UseVectorEyeColorInImages(true)
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			die := r1.Intn(6) + 1
			dieImage := fmt.Sprintf("images/dice/%d.png", die)
			dieImage = sdk_wrapper.GetDataPath(dieImage)

			images.DisplayAnimatedGif(sdk_wrapper.GetDataPath("images/dice/roll-the-dice.gif"), images.ANIMATED_GIF_SPEED_FASTEST, 1, false)
			images.DisplayImage(dieImage, 100, false)
			audio.PlaySound(audio.SYSTEMSOUND_WIN)
			voice.SayText(fmt.Sprintf("You rolled a %d", die))
			images.DisplayImageWithTransition(dieImage, 1000, images.IMAGE_TRANSITION_FADE_OUT, 10)

			stop <- true
			return
		}
	}
}
