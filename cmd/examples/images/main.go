package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/images"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/motors"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/voice"
	"image/color"
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
			for l := 1; l < 64; l++ {
				images.WriteText("TEXT", float64(l), true, 30, true)
			}
			voice.SayText("Text!")
			images.WriteText("TEXT", 64, true, 1000, false)
			images.WriteColoredText("This", 32, true, color.RGBA{255, 0, 0, 255}, 500, true)
			images.WriteColoredText("is", 32, true, color.RGBA{255, 255, 0, 255}, 500, true)
			images.WriteColoredText("some", 32, true, color.RGBA{0, 255, 0, 255}, 500, true)
			images.WriteColoredText("text", 32, true, color.RGBA{0, 255, 255, 255}, 500, true)
			voice.SayText("Image display")
			images.DisplayImage("data/images/birthday-cake.jpg", 5000, true)
			voice.SayText("Image transitions")
			images.WriteText("SLIDE", 32, true, 2000, true)
			images.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, images.IMAGE_TRANSITION_SLIDE_LEFT, 10)
			images.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, images.IMAGE_TRANSITION_SLIDE_RIGHT, 10)
			images.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, images.IMAGE_TRANSITION_SLIDE_DOWN, 10)
			images.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, images.IMAGE_TRANSITION_SLIDE_UP, 10)
			images.WriteText("FADE", 32, true, 2000, true)
			images.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, images.IMAGE_TRANSITION_FADE_IN, 10)
			images.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, images.IMAGE_TRANSITION_FADE_OUT, 10)
			voice.SayText("Animated gifs")

			images.SetBackgroundColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
			images.DisplayAnimatedGif("data/images/animated.gif", images.ANIMATED_GIF_SPEED_SLOW, 3, false)
			images.DisplayAnimatedGif("data/images/animated2.gif", images.ANIMATED_GIF_SPEED_NORMAL, 3, false)
			images.DisplayAnimatedGif("data/images/animated3.gif", images.ANIMATED_GIF_SPEED_FAST, 3, false)
			images.WriteText("ENJOY!", 32, true, 5000, true)
			stop <- true
			return
		}
	}
}
