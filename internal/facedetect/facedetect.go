package facedetect

import (
	"fmt"
	"gocv.io/x/gocv"
	"image/color"
)

func Some() {
	img := gocv.IMRead("/path/to/your/photo.jpg", gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("Error reading image")
		return
	}
	defer img.Close()

	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()
	if !classifier.Load("/home/anwarzadeh/Desktop/face-detection/cmd/haarcascade_frontalface_default.xml") {
		fmt.Println("Error loading cascade classifier.")
		return
	}

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	faces := classifier.DetectMultiScale(gray)
	if len(faces) == 0 {
		fmt.Println("No faces detected")
	} else {

		for _, face := range faces {
			gocv.Rectangle(&img, face, color.RGBA{0, 255, 0, 0}, 2)
		}
	}

	window := gocv.NewWindow("Face Detection")
	defer window.Close()
	window.IMShow(img)
	window.WaitKey(0)
}
