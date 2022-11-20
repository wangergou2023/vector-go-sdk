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

/*
import anki_vector
import sys

def main():
    sentences = [
            "Fortune favors the bold",
            "I think, therefore I am",
            "Time is money",
            "I came, I saw, I conquered",
            "When life gives you lemons, make lemonade",
            "Practice makes perfect",
            "Knowledge is power",
            "Have no fear of perfection, you'll never reach it",
            "No pain no gain",
            "That which does not kill us makes us stronger"
            ]

    with anki_vector.Robot(cache_animation_lists=False) as robot:
        phrase = ""
        if len(sys.argv)>0:
            phrase = sys.argv[1]

        if len(phrase)==0:
            from random import randrange
            i = randrange(10)
            phrase = sentences[i]

        print("Say " + phrase)
        robot.behavior.say_text(phrase, False)

if __name__ == "__main__":
    main()

*/

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
