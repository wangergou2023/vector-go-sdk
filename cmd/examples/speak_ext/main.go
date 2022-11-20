package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"log"
	"math/rand"
	"time"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	var sentence = flag.String("sentence", "", "Sentence to say")
	flag.Parse()

	v, err := vector.NewEP(*serial)
	if err != nil {
		log.Fatal(err)
	}

	sentences := [10]string{
		"Fortune favors the bold",
		"I think, therefore I am",
		"Time is money",
		"I came, I saw, I conquered",
		"When life gives you lemons, make lemonade",
		"Practice makes perfect",
		"Knowledge is power",
		"Have no fear of perfection, you'll never reach it",
		"No pain no gain",
		"That which does not kill us makes us stronger",
	}

	var phrase = ""

	if *sentence == "" {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		phrase = sentences[r1.Intn(10)]
	} else {
		phrase = *sentence
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
			_, _ = v.Conn.SayText(
				ctx,
				&vectorpb.SayTextRequest{
					Text:           phrase,
					UseVectorVoice: true,
					DurationScalar: 1.0,
				},
			)
			stop <- true
			return
		}
	}
}
