package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/fforchino/vector-go-sdk/pkg/oskrpb"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	WIFI_MAX       = 100
	SEARCH_SPEED   = 50
	AVOID_SPEED    = 25
	MIN_PROXIMITY  = 10
	CHECK_INTERVAL = 100 * time.Millisecond
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	sdk_wrapper.InitSDK(*serial)
	conn, err := grpc.Dial(sdk_wrapper.Robot.Cfg.Target+":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := oskrpb.NewOSKRServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			Navigate(ctx, client)
			stop <- true
			return
		}
	}
}

func Navigate(ctx context.Context, client oskrpb.OSKRServiceClient) {
	maxSignal := 0
	prevSignal := 0
	for {
		proximity := GetProximitySensorData()

		// Get wifi signal strength using grpc
		wifiSignalRes, err := client.GetWifiSignalStrength(ctx, &oskrpb.WifiSignalStrengthRequest{})
		if err != nil {
			log.Fatalf("could not get signal strength: %v", err)
		}
		wifiSignal := int(wifiSignalRes.GetSignalStrength())
		log.Printf("Wifi signal strength: %d", wifiSignal)

		if wifiSignal > maxSignal {
			maxSignal = wifiSignal
		}

		if proximity <= MIN_PROXIMITY {
			MoveRobot(-AVOID_SPEED, AVOID_SPEED)
		} else if wifiSignal < prevSignal {
			MoveRobot(-SEARCH_SPEED, SEARCH_SPEED)
		} else {
			MoveRobot(SEARCH_SPEED, SEARCH_SPEED)
		}

		if wifiSignal == WIFI_MAX {
			fmt.Println("Maximum signal strength reached.")
			MoveRobot(0, 0)
			break
		}

		prevSignal = wifiSignal
		time.Sleep(CHECK_INTERVAL)
	}
}

func MoveRobot(leftSpeed int, rightSpeed int) {
	// Placeholder for the API call
	//sdk_wrapper.DriveWheelsForward(float32(leftSpeed), float32(rightSpeed), 1, 1)
}

func GetProximitySensorData() int {
	// Placeholder for the API call
	return 0
}
