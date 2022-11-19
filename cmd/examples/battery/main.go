package main

import (
	"context"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"log"
)

func main() {
	v, err := vector.NewEP()
	if err != nil {
		log.Fatal(err)
	}

	bs, err := v.Conn.BatteryState(
		context.Background(),
		&vectorpb.BatteryStateRequest{},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bs)
}
