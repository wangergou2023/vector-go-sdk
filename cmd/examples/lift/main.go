package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"log"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	v, err := vector.NewEP(*serial)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = v.BehaviorControl(ctx, start, stop, nil)
	}()

	for {
		select {
		case <-start:
			_, err := v.Conn.SetLiftHeight(
				ctx,
				&vectorpb.SetLiftHeightRequest{
					HeightMm:          250,
					MaxSpeedRadPerSec: 1,
					IdTag:             2000001,
					DurationSec:       .001,
				},
			)
			if err != nil {
				fmt.Println(err)
			}
			stop <- true
			return
		}
	}

}
