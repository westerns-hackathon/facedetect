package handler

import (
	"face-detection/internal/facedetect"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) PostPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Не удалось распарсить форму", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Сохранение загруженного файла
	fileName := fmt.Sprintf("photo_%d.jpg", time.Now().Unix())
	inputPath := "storage/photo/" + fileName
	outFile, err := os.Create(inputPath)
	if err != nil {
		http.Error(w, "Ошибка при создании файла на сервере", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	outputDir := "storage/detected_photo/"
	if err := facedetect.DetectFaces(inputPath, outputDir); err != nil {
		http.Error(w, "Ошибка при распознавании лиц", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Фото успешно загружено и обработано!")
}
