package sdk_wrapper

import (
	"bytes"
	"fmt"
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

// Saves current image at high resolution (1280 x 720) on a file in jpg format
// This doesn't seem to work on my production Vector on Wirepod. Probably vector-cloud on the robot needs to be
// updated. But since it's a production robot I'm stuck to 1.8...
// Anyways the image is saved at the regular 360p size.

func SaveHiResCameraPicture(fileName string) error {
	i, err := Robot.Conn.CaptureSingleImage(
		ctx,
		&vectorpb.CaptureSingleImageRequest{
			EnableHighResolution: true,
		},
	)
	if err == nil {
		var imageData = i.GetData()
		var mat gocv.Mat
		mat, err = gocv.IMDecode(imageData, -1)
		var size = mat.Size()
		println(fmt.Sprintf("\nGot an image %dx%d\n", size[1], size[0]))
		if err == nil && !mat.Empty() {
			buf, err2 := gocv.IMEncode(".jpg", mat)
			if err2 == nil {
				err = os.WriteFile(fileName, buf.GetBytes(), 0644)
			}
			err = err2
		}
	}

	return err
}

func GetCameraPictureOpenCv() gocv.Mat {
	var mat gocv.Mat
	i, err := Robot.Conn.CaptureSingleImage(
		ctx,
		&vectorpb.CaptureSingleImageRequest{
			EnableHighResolution: true,
		},
	)
	if err == nil {
		var imageData = i.GetData()
		var mat2 gocv.Mat
		mat2, err = gocv.IMDecode(imageData, -1)
		if err == nil {
			mat = mat2
		}
	}
	return mat
}
