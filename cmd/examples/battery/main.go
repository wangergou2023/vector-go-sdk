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

	bs, err := v.Conn.BatteryState(
		context.Background(),
		&vectorpb.BatteryStateRequest{},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bs)
}
