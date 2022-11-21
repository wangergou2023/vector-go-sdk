package sdk_wrapper

import (
	"bytes"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"gocv.io/x/gocv"
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

// Enables camera, saves current image (1280 x 720) on a file in jpg format, disables camera
func SaveHiResCameraPicture(fileName string) error {
	if !camStreamEnable {
		EnableCameraStream()
	}
	i, err := Robot.Conn.CaptureSingleImage(
		ctx,
		&vectorpb.CaptureSingleImageRequest{},
	)
	if err != nil {
		var imageData = i.GetData()
		var mat gocv.Mat
		mat, err = gocv.IMDecode(imageData, -1)
		if err == nil && !mat.Empty() {
			buf, err2 := gocv.IMEncode("jpg", mat)
			if err2 == nil {
				err = os.WriteFile(fileName, buf.GetBytes(), 0644)
			}
			err = err2
		}
	}
	DisableCameraStream()
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
