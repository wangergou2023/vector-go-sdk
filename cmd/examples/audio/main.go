package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"log"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	var sentence = flag.String("sentence", "", "Sentence to say")
	var audioFile = flag.String("audiofile", "", "Audio file to stream")
	var volume = flag.String("volume", "", "Volume (0-100)")
	flag.Parse()

	v, err := vector.NewEP(*serial)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = v.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			if *sentence != "" {
				_, _ = v.Conn.SayText(
					ctx,
					&vectorpb.SayTextRequest{
						Text:           *sentence,
						UseVectorVoice: true,
						DurationScalar: 1.0,
					},
				)
			}
			if *audioFile != "" {
				_, _ = v.Conn.ExternalAudioStreamPlayback(
					ctx,
					&vectorpb.ExternalAudioStreamRequest{
						AudioRequestType: ExternalAudioStreamRequest_AudioStreamPrepare,
					},
				)
			}
			stop <- true
			return
		}
	}
}
