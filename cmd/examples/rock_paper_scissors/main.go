package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	sdk_wrapper "github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

/*
 Rock paper scissors icons found on flaticon.com

 <a href="https://www.flaticon.com/free-icons/rock-paper-scissors" title="rock paper scissors icons">Rock paper scissors icons created by Freepik - Flaticon</a>
*/

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
			playGame(3)
			stop <- true
			return
		}
	}
}

func playGame(numSteps int) {
	options := [3]string{
		"rock",
		"paper",
		"scissors",
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i <= numSteps; i++ {
		sdk_wrapper.SayText("one, two, three!")

		myMove := options[r1.Intn(len(options))]
		sdk_wrapper.DisplayImage("data/images/"+myMove+".png", 5000, false)
		sdk_wrapper.PlaySound("data/sounds/quick-win.wav", 100)
		sdk_wrapper.SayText(myMove + "!")

		fName := "/tmp/rps.jpg"
		err := sdk_wrapper.SaveHiResCameraPicture(fName)
		if err == nil {
			if err == nil {
				cmd := exec.Command("python", "hand_detection.py", fName)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("FRAME %d, Output: %s\n", i, out.String())

				win := 0
				answer := ""
				userMove := ""
				if out.String() == "0" {
					// User plays "rock"
					userMove = "rock"
					if myMove == "paper" {
						win = 1
					} else if myMove == "scissors" {
						win = -1
					}
				} else if out.String() == "2" {
					// User plays "scissors"
					userMove = "scissors"
					if myMove == "rock" {
						win = 1
					} else if myMove == "paper" {
						win = -1
					}
				} else if out.String() == "5" {
					// User plays "paper"
					userMove = "paper"
					if myMove == "scissors" {
						win = 1
					} else if myMove == "rock" {
						win = -1
					}
				} else {
					answer = "Sorry, I don't get it"
				}

				if answer == "" {
					answer = "You put " + userMove + ", I put " + myMove + ". "

					switch win {
					case -1:
						answer = answer + "You win!"
						break
					case 1:
						answer = answer + "I win!"
						break
					default:
						answer = answer + "It's a draw!"
						break
					}
				}
				sdk_wrapper.SayText(answer)
				time.Sleep(time.Duration(5000) * time.Millisecond)
			} else {
				println("OPENCV Python script not found!")
			}
		}
	}
}
