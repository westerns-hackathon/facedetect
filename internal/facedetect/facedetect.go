package facedetect

import (
	"face-detection/internal/db"
	"face-detection/internal/model"
	"fmt"
	"github.com/Kagami/go-face"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
)

const modelsDir = "internal/facedetect/models"

func getModelPaths() (string, string) {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	modelPath := filepath.Join(cwd, "cmd", "res10_300x300_ssd_iter_140000.caffemodel")
	prototxtPath := filepath.Join(cwd, "cmd", "deploy.prototxt")

	return modelPath, prototxtPath
}

func DetectFaces(inputPath, outputDir string) error {
	img := gocv.IMRead(inputPath, gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("ошибка при чтении изображения")
	}
	defer img.Close()

	modelPath, prototxtPath := getModelPaths()

	net := gocv.ReadNet(modelPath, prototxtPath)

	defer net.Close()

	if net.Empty() {
		return fmt.Errorf("ошибка при загрузке DNN модели")
	}

	blob := gocv.BlobFromImage(img, 1.0, image.Point{300, 300}, gocv.NewScalar(104, 177, 123, 0), false, false)
	defer blob.Close()
	net.SetInput(blob, "")
	detections := net.Forward("")

	if detections.Empty() {
		return fmt.Errorf("не удалось получить детекции")
	}

	detectionMat := detections.Reshape(1, detections.Total()/7)
	numDetections := detectionMat.Rows()

	for i := 0; i < numDetections; i++ {
		confidence := detectionMat.GetFloatAt(i, 2)
		if confidence > 0.8 {
			left := int(detectionMat.GetFloatAt(i, 3) * float32(img.Cols()))
			top := int(detectionMat.GetFloatAt(i, 4) * float32(img.Rows()))
			right := int(detectionMat.GetFloatAt(i, 5) * float32(img.Cols()))
			bottom := int(detectionMat.GetFloatAt(i, 6) * float32(img.Rows()))

			rect := image.Rect(left, top, right, bottom)

			if right-left > 50 && bottom-top > 50 {
				gocv.Rectangle(&img, rect, color.RGBA{21, 104, 100, 0}, 1)
				gocv.PutText(&img, fmt.Sprintf("Confidence: %.2f", confidence), image.Pt(left, top-10),
					gocv.FontHersheySimplex, 0.5, color.RGBA{255, 0, 0, 0}, 1)
			}
		}
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("не удалось создать директорию: %v", err)
	}
	outputPath := filepath.Join(outputDir, "detected_"+filepath.Base(inputPath))
	if ok := gocv.IMWrite(outputPath, img); !ok {
		return fmt.Errorf("не удалось сохранить изображение")
	}
	fmt.Printf("Фото с лицами сохранено в %s\n", outputPath)
	return nil
}

func MatchFaces(firstImage, secondImage string) (string, error) {

	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		return "", fmt.Errorf("не удалось инициализировать распознаватель лиц: %v", err)
	}
	defer rec.Close()

	// Получение лиц на обоих изображениях
	faces1, err := rec.RecognizeFile(firstImage)
	if err != nil || len(faces1) == 0 {
		return "", fmt.Errorf("не удалось распознать лицо на первом изображении: %v", err)
	}

	faces2, err := rec.RecognizeFile(secondImage)
	if err != nil || len(faces2) == 0 {
		return "", fmt.Errorf("не удалось распознать лицо на втором изображении: %v", err)
	}

	// Сравнение дескрипторов
	distance := face.SquaredEuclideanDistance(faces1[0].Descriptor, faces2[0].Descriptor)
	resultMessage := fmt.Sprintf("Расстояние между лицами: %.2f", distance)
	if distance < 0.6 {
		return "Лица совпадают. " + resultMessage, nil
	}

	return "Лица не совпадают. " + resultMessage, nil
}

func GetFaceDescriptors(outputDir string) ([][]float32, error) {
	var descriptors [][]float32

	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать распознаватель лиц: %v", err)
	}
	defer rec.Close()

	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		log.Printf("Обрабатываем файл: %s", path)

		faces, err := rec.RecognizeFile(path)
		if err != nil {
			return fmt.Errorf("не удалось распознать лица на изображении %s: %v", path, err)
		}

		for _, f := range faces {
			descriptors = append(descriptors, f.Descriptor[:]) // Берем срез от Descriptor
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка при обработке изображений: %v", err)
	}

	return descriptors, nil
}

func CosineSimilarity(a, b []float32) float32 {
	var dotProduct float32
	var normA float32
	var normB float32

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
func FindMatchingFaces(descriptor []float32, db db.Storage) ([]model.Face, error) {
	faces, err := db.GetAllFaces()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех лиц: %v", err)
	}

	var matchingFaces []model.Face
	for _, face := range faces {
		similarity := CosineSimilarity(descriptor, face.Descriptor)
		if similarity > 0.9 {
			matchingFaces = append(matchingFaces, face)
		}
	}
	return matchingFaces, nil
}
