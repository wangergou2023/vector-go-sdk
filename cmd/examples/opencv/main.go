package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"math"
)

// What it does:
//
// This example detects how many fingers you hold up in front of the camera.
//

const MinimumArea = 3000

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	v, err := vector.NewEP(*serial)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = v.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			sdk_wrapper.PlayAnimation("anim_generic_look_up_01", 0, false, false, false)
			doOpenCVStuff()
			stop <- true
			return
		}
	}
}

func doOpenCVStuff() {
	img := gocv.NewMat()
	defer img.Close()

	imgGrey := gocv.NewMat()
	defer imgGrey.Close()

	imgBlur := gocv.NewMat()
	defer imgBlur.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	hull := gocv.NewMat()
	defer hull.Close()

	defects := gocv.NewMat()
	defer defects.Close()

	green := color.RGBA{0, 255, 0, 0}

	fmt.Printf("Start reading camera")

	sdk_wrapper.EnableCameraStream()
	for {
		img = sdk_wrapper.GetCameraPictureOpenCv()

		if img.Empty() {
			continue
		}

		// cleaning up image
		gocv.CvtColor(img, &imgGrey, gocv.ColorBGRToGray)
		gocv.GaussianBlur(imgGrey, &imgBlur, image.Pt(35, 35), 0, 0, gocv.BorderDefault)
		gocv.Threshold(imgBlur, &imgThresh, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)

		// now find biggest contour
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		c := getBiggestContour(contours)

		gocv.ConvexHull(c, &hull, true, false)
		gocv.ConvexityDefects(c, hull, &defects)

		var angle float64
		defectCount := 0
		for i := 0; i < defects.Rows(); i++ {
			start := c.At(int(defects.GetIntAt(i, 0)))
			end := c.At(int(defects.GetIntAt(i, 1)))
			far := c.At(int(defects.GetIntAt(i, 2)))

			a := math.Sqrt(math.Pow(float64(end.X-start.X), 2) + math.Pow(float64(end.Y-start.Y), 2))
			b := math.Sqrt(math.Pow(float64(far.X-start.X), 2) + math.Pow(float64(far.Y-start.Y), 2))
			c := math.Sqrt(math.Pow(float64(end.X-far.X), 2) + math.Pow(float64(end.Y-far.Y), 2))

			// apply cosine rule here
			angle = math.Acos((math.Pow(b, 2)+math.Pow(c, 2)-math.Pow(a, 2))/(2*b*c)) * 57

			// ignore angles > 90 and highlight rest with dots
			if angle <= 90 {
				defectCount++
				gocv.Circle(&img, far, 1, green, 2)
			}
		}

		status := fmt.Sprintf("defectCount: %d", defectCount+1)
		println(status)
	}
}

func getBiggestContour(contours gocv.PointsVector) gocv.PointVector {
	var area float64
	index := 0
	for i := 0; i < contours.Size(); i++ {
		newArea := gocv.ContourArea(contours.At(i))
		if newArea > area {
			area = newArea
			index = i
		}
	}
	return contours.At(index)
}
