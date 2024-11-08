package facedetect

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
	"path/filepath"
)

func DetectFaces(inputPath, outputDir string) error {
	// Read the input image
	img := gocv.IMRead(inputPath, gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("ошибка при чтении изображения")
	}
	defer img.Close()

	// Load the pre-trained DNN face detection model
	net := gocv.ReadNet(
		"/home/anwarzadeh/Desktop/face-detection/cmd/res10_300x300_ssd_iter_140000.caffemodel",
		"/home/anwarzadeh/Desktop/face-detection/cmd/deploy.prototxt",
	)

	net.SetPreferableBackend(gocv.NetBackendDefault)
	net.SetPreferableTarget(gocv.NetTargetCPU)

	if net.Empty() {
		return fmt.Errorf("ошибка при загрузке DNN модели")
	}
	defer net.Close()

	blob := gocv.BlobFromImage(img, 1.0, image.Point{300, 300}, gocv.NewScalar(104, 177, 123, 0), false, false)
	defer blob.Close()

	net.SetInput(blob, "")
	detections := net.Forward("")

	if detections.Empty() {
		return fmt.Errorf("не удалось получить детекции")
	}

	detectionMat := detections.Reshape(1, detections.Total()/7)

	numDetections := detectionMat.Rows()

	fmt.Printf("Number of detections: %d\n", numDetections)

	for i := 0; i < numDetections; i++ {
		confidence := detectionMat.GetFloatAt(i, 2)

		fmt.Printf("Detection %d, Confidence: %f, Coordinates: (%f, %f, %f, %f)\n",
			i, confidence, detectionMat.GetFloatAt(i, 3), detectionMat.GetFloatAt(i, 4),
			detectionMat.GetFloatAt(i, 5), detectionMat.GetFloatAt(i, 6))

		if confidence > 0.5 {

			left := int(detectionMat.GetFloatAt(i, 3) * float32(img.Cols()))
			top := int(detectionMat.GetFloatAt(i, 4) * float32(img.Rows()))
			right := int(detectionMat.GetFloatAt(i, 5) * float32(img.Cols()))
			bottom := int(detectionMat.GetFloatAt(i, 6) * float32(img.Rows()))

			rect := image.Rect(left, top, right, bottom)
			gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)
		}
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("не удалось создать директорию: %v", err)
	}

	// Save the image with the detected faces to the output directory
	outputPath := filepath.Join(outputDir, "detected_"+filepath.Base(inputPath))
	if ok := gocv.IMWrite(outputPath, img); !ok {
		return fmt.Errorf("не удалось сохранить изображение")
	}
	fmt.Printf("Фото с лицами сохранено в %s\n", outputPath)

	return nil
}
