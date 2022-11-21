package sdk_wrapper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"os"
	"time"
)

const IMAGE_TRANSITION_NONE = 0
const IMAGE_TRANSITION_FADE_IN = 1
const IMAGE_TRANSITION_SLIDE_LEFT = 2
const IMAGE_TRANSITION_SLIDE_RIGHT = 3
const IMAGE_TRANSITION_SLIDE_UP = 4
const IMAGE_TRANSITION_SLIDE_DOWN = 5

const TRANSITION_SLICES = 20

func TextOnImg(text string, size float64, isBold bool) []byte {
	bgImage := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 184, Y: 96},
	})
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)

	var fontName = "DroidSans"
	if isBold {
		fontName = fontName + "-Bold"
	}

	if err := dc.LoadFontFace("data/fonts/"+fontName+".ttf", size); err != nil {
		fmt.Println(err)
		return nil
	}

	x := float64(imgWidth / 2)
	y := float64((imgHeight / 2))
	maxWidth := float64(imgWidth) - 35.0
	dc.SetColor(color.RGBA{0, 255, 0, 255}) // Green
	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)
	buf := new(bytes.Buffer)
	bitmap := convertPixelsToRawBitmap(dc.Image())
	for _, ui := range bitmap {
		binary.Write(buf, binary.LittleEndian, ui)
	}
	os.WriteFile("/tmp/test.raw", buf.Bytes(), 0644)
	return buf.Bytes()
}

func DataOnImg(fileName string) []byte {
	inFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer inFile.Close()

	src, _, err := image.Decode(inFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bgImage := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 184, Y: 96},
	})
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)

	var dst = src
	if src.Bounds().Dy() > src.Bounds().Dx() {
		dst = resize.Resize(0, uint(imgHeight), src, resize.Bicubic)
	} else {
		dst = resize.Resize(uint(imgWidth), 0, src, resize.Bicubic)
	}
	dc.DrawImage(dst, (imgWidth-dst.Bounds().Dx())/2, (imgHeight-dst.Bounds().Dy())/2)

	buf := new(bytes.Buffer)
	bitmap := convertPixelsToRawBitmap(dc.Image())
	for _, ui := range bitmap {
		binary.Write(buf, binary.LittleEndian, ui)
	}
	os.WriteFile("/tmp/test.raw", buf.Bytes(), 0644)
	return buf.Bytes()
}

func DataOnImgWithTransition(fileName string, transition int, pct int) []byte {
	inFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer inFile.Close()

	src, _, err := image.Decode(inFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bgImage := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 184, Y: 96},
	})
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)

	var dst = src
	if src.Bounds().Dy() > src.Bounds().Dx() {
		dst = resize.Resize(0, uint(imgHeight), src, resize.Bilinear)
	} else {
		dst = resize.Resize(uint(imgWidth), 0, src, resize.Bilinear)
	}
	x := (imgWidth - dst.Bounds().Dx()) / 2
	y := (imgHeight - dst.Bounds().Dy()) / 2
	if transition == IMAGE_TRANSITION_SLIDE_RIGHT {
		xStart := -1 * dst.Bounds().Dx()
		xEnd := (imgWidth - dst.Bounds().Dx()) / 2
		x = xStart + (xEnd-xStart)*pct/100
		//println(x)
	} else if transition == IMAGE_TRANSITION_SLIDE_LEFT {
		xStart := imgWidth + dst.Bounds().Dx()
		xEnd := (imgWidth - dst.Bounds().Dx()) / 2
		x = xStart - (xStart-xEnd)*pct/100
		//println(x)
	} else if transition == IMAGE_TRANSITION_SLIDE_DOWN {
		yStart := -1 * dst.Bounds().Dy()
		yEnd := (imgHeight - dst.Bounds().Dy()) / 2
		y = yStart + (yEnd-yStart)*pct/100
		//println(y)
	} else if transition == IMAGE_TRANSITION_SLIDE_UP {
		yStart := imgHeight + dst.Bounds().Dy()
		yEnd := (imgHeight - dst.Bounds().Dy()) / 2
		y = yStart - (yStart-yEnd)*pct/100
		//println(y)
	}
	dc.DrawImage(dst, x, y)

	buf := new(bytes.Buffer)
	bitmap := convertPixelsToRawBitmap(dc.Image())
	for _, ui := range bitmap {
		binary.Write(buf, binary.LittleEndian, ui)
	}
	os.WriteFile("/tmp/test.raw", buf.Bytes(), 0644)
	return buf.Bytes()
}

func WriteText(text string, size float64, isBold bool, duration int, blocking bool) {
	faceBytes := TextOnImg(text, size, isBold)
	displayFaceImage(faceBytes, duration, blocking)
}

func DisplayImage(imageFile string, duration int, blocking bool) {
	faceBytes := DataOnImg(imageFile)
	displayFaceImage(faceBytes, duration, blocking)
}

func DisplayImageWithTransition(imageFile string, duration int, transition int, numSteps int) {
	if numSteps == 0 || transition == IMAGE_TRANSITION_NONE {
		faceBytes := DataOnImg(imageFile)
		displayFaceImage(faceBytes, duration, true)
	} else {
		for i := 0; i <= numSteps; i++ {
			pctProgress := i * 100 / numSteps
			tmpFaceBytes := DataOnImgWithTransition(imageFile, transition, pctProgress)
			displayFaceImage(tmpFaceBytes, 100, false)
		}
	}
}

func displayFaceImage(faceBytes []byte, duration int, blocking bool) {
	_, _ = Robot.Conn.DisplayFaceImageRGB(
		ctx,
		&vectorpb.DisplayFaceImageRGBRequest{
			FaceData:         faceBytes,
			DurationMs:       uint32(duration),
			InterruptRunning: true,
		},
	)
	if blocking {
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}
}

func convertPixesTo16BitRGB(r uint32, g uint32, b uint32, a uint32) uint16 {
	// RGB565: 5 bits for red, 6 for green, 5 for blue
	R, G, B := int8(r/257), int8(g/257), int8(b/257)

	bx := (B >> 3)
	gx := (G >> 2) << 5
	rx := (R >> 3) << 11

	return uint16(rx | gx | bx)
}

func convertPixelsToRawBitmap(image image.Image) []uint16 {
	imgHeight, imgWidth := image.Bounds().Max.Y, image.Bounds().Max.X
	bitmap := make([]uint16, imgWidth*imgHeight)

	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			bitmap[(y)*imgWidth+(x)] = convertPixesTo16BitRGB(image.At(x, y).RGBA())
		}
	}
	return bitmap
}
