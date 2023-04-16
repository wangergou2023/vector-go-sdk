package main

import (
	"context"
	"flag"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"os/exec"
	"regexp"
	"strconv"
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

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			Navigate()
			stop <- true
			return
		}
	}
}

func Navigate() {
	maxSignal := 0
	prevSignal := 0
	for {
		proximity := GetProximitySensorData()
		wifiSignal := GetWifiSignalStrength()

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
}

func GetProximitySensorData() int {
	// Placeholder for the API call
	return 0
}

func GetWifiSignalStrength() int {
	cmd := exec.Command("sh", "-c", "iwconfig wlan0 | grep -i --color signal")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return 0
	}

	signalRegexp := regexp.MustCompile(`Signal level=(-?\d+) dBm`)
	matches := signalRegexp.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		fmt.Println("Error parsing signal level")
		return 0
	}

	signal, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Println("Error converting signal level to integer:", err)
		return 0
	}

	// Convert dBm to a percentage value between 0 and WIFI_MAX
	percentage := int((float64(signal+100) / 70) * float64(WIFI_MAX))
	if percentage < 0 {
		percentage = 0
	} else if percentage > WIFI_MAX {
		percentage = WIFI_MAX
	}

	return percentage
}
