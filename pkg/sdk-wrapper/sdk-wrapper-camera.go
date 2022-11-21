package sdk_wrapper

import (
	"bytes"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
)

var camStreamEnable bool = false
var camStreamClient vectorpb.ExternalInterface_CameraFeedClient

func EnableCameraStream() {
	camStreamClient, _ = Robot.Conn.CameraFeed(ctx, &vectorpb.CameraFeedRequest{})
	camStreamEnable = true
}

func DisableCameraStream() {
	camStreamEnable = false
	camStreamClient = nil
}

func ProcessCameraStream() image.Image {
	i := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 640, Y: 360},
	})
	if camStreamEnable {
		response, _ := camStreamClient.Recv()
		imageBytes := response.GetData()
		img, _, _ := image.Decode(bytes.NewReader(imageBytes))
		return img
	} else {
		for j := range i.Pix {
			i.Pix[j] = uint8(rand.Uint32())
		}
	}

	return i.SubImage(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 640, Y: 360},
	})
}

// Enables camera, saves current image on a file in jpg format, disables camera
func GetCameraPicture() image.Image {
	if !camStreamEnable {
		EnableCameraStream()
	}
	var img = ProcessCameraStream()
	DisableCameraStream()
	return img
}

// Enables camera, saves current image on a file in jpg format, disables camera
func SaveCameraPicture(fileName string) error {
	var img = GetCameraPicture()
	f, err := os.Create(fileName)
	if err == nil {
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(f, img, &opt)
	}
	return err
}

/*
func ProcessCameraStream() image {
	camStream := mjpeg.NewStream()
	i := image.NewGray(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 640, Y: 360},
	})
	go func() {
		for {
			if camStreamEnable {
				response, _ := camStreamClient.Recv()
				imageBytes := response.GetData()
				img, _, _ := image.Decode(bytes.NewReader(imageBytes))
				camStream.Update(img)
			} else {
				for j := range i.Pix {
					i.Pix[j] = uint8(rand.Uint32())
				}
				time.Sleep(time.Second)
				camStream.Update(i)
			}
		}
	}()
}

*/
