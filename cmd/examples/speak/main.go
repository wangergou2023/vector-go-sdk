package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
)

func main() {
	var guid = "q0MulUycT6aThkfFavXoog=="
	/*
		var token = get_session_token("filippo.forchino@gmail.com", "Suka99!!!")
		st := token["session"].(map[string]interface{})
		stt := fmt.Sprintf("%s", st["session_token"])
		print("Token: " + stt)
		v, err := vector.New(
			vector.WithTarget(os.Getenv("BOT_TARGET")),
			//vector.WithToken(os.Getenv("BOT_TOKEN")),
			vector.WithToken(stt),
		)
	*/
	v, err := vector.New(
		vector.WithTarget(os.Getenv("BOT_TARGET")),
		//vector.WithToken(os.Getenv("BOT_TOKEN")),
		vector.WithToken(guid),
	)

	print("OK")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)
	print("OK")

	go func() {
		print("OK")
		_ = v.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			print("START")
			_, _ = v.Conn.SayText(
				ctx,
				&vectorpb.SayTextRequest{
					Text:           "hello, hello, hello",
					UseVectorVoice: true,
					DurationScalar: 1.0,
				},
			)
			stop <- true
			return
		}
	}
}

func get_session_token(username string, password string) map[string]interface{} {
	ret := ""
	api := "https://accounts.api.anki.com/1/sessions"

	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)

	req, err := http.NewRequest("POST", api, strings.NewReader(params.Encode()))
	req.Header.Set("User-Agent", "Vector-sdk/0.7.2")
	req.Header.Set("Anki-App-Key", "aung2ieCho3aiph7Een3Ei")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Get Session Token Response Status:\n", resp.Status)
	if resp.Status != "200 OK" {
		panic("Authentication error")
	} else {
		body, _ := io.ReadAll(resp.Body)
		ret = string(body)
	}
	print(ret)
	var jsonObj map[string]interface{}
	json.Unmarshal([]byte(ret), &jsonObj)
	return jsonObj
}
