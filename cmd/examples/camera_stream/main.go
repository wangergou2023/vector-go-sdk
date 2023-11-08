package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sync"
	"time"

	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	var port = flag.Int("port", 80, "listen port")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)

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
			sdk_wrapper.MoveHead(3.0)
			sdk_wrapper.EnableCameraStream()

			m := NewMJPEGStream()
			http.Handle("/", m)

			go func() {
				log.Println("Starting server on", addr)
				err := http.ListenAndServe(addr, nil)
				if err != nil {
					log.Println(err)
				}
			}()

			for {
				img := sdk_wrapper.ProcessCameraStream()
				if img != nil {
					m.UpdateJPEG(imageToJPEG(img))
				}
			}

			stop <- true
			return
		}
	}
}

// MJPEGStream is a structure that allows to serve a stream of JPEG images as MJPEG over HTTP.
type MJPEGStream struct {
	boundary string
	frame    []byte
	lock     sync.Mutex
}

// NewMJPEGStream returns a new MJPEGStream object.
func NewMJPEGStream() *MJPEGStream {
	return &MJPEGStream{
		boundary: "--BOUNDARY",
		frame:    make([]byte, 0),
	}
}

// UpdateJPEG updates the current frame with a new JPEG image.
func (s *MJPEGStream) UpdateJPEG(jpeg []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.frame = jpeg
}

// ServeHTTP handles HTTP requests to the stream.
func (s *MJPEGStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving MJPEG stream.")

	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+s.boundary)
	multipartWriter := multipart.NewWriter(w)
	multipartWriter.SetBoundary(s.boundary)
	for {
		s.lock.Lock()
		frame := s.frame
		s.lock.Unlock()
		if len(frame) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		imageWriter, err := multipartWriter.CreatePart(textproto.MIMEHeader{
			"Content-type":   []string{"image/jpeg"},
			"Content-length": []string{fmt.Sprint(len(frame))},
		})
		if err != nil {
			log.Println(err)
			break
		}
		imageWriter.Write(frame)
	}
}

// imageToJPEG converts an image.Image to a JPEG byte slice.
func imageToJPEG(img image.Image) []byte {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	return buf.Bytes()
}
