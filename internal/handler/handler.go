package handler

import (
	"encoding/json"
	"face-detection/internal/db"
	"face-detection/internal/facedetect"
	"face-detection/internal/model"
	"fmt"
	"github.com/Kagami/go-face"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Handler struct {
	Recognizer *face.Recognizer
	storage    db.Storage
}

func NewHandler(data *face.Recognizer, storage db.Storage) *Handler {
	return &Handler{data, storage}
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

	log.Printf("Путь к загруженному изображению для распознавания лиц: %s\n", inputPath)

	outputDir := "storage/detected_photo"
	if err := facedetect.DetectFaces(inputPath, outputDir); err != nil {
		http.Error(w, "Ошибка при обработке изображения", http.StatusInternalServerError)
		return
	}

	outputPath := filepath.Join(outputDir, "detected_"+fileName)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		http.Error(w, "Обработанный файл не найден", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"image_path": outputPath,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}

func (h *Handler) PostFaceMatchHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Не удалось распарсить форму", http.StatusBadRequest)
		return
	}

	file1, _, err := r.FormFile("first_image")
	if err != nil {
		http.Error(w, "Ошибка при получении первого изображения", http.StatusBadRequest)
		return
	}
	defer file1.Close()

	file2, _, err := r.FormFile("second_image")
	if err != nil {
		http.Error(w, "Ошибка при получении второго изображения", http.StatusBadRequest)
		return
	}
	defer file2.Close()

	uploadsDir := filepath.Join("storage", "uploads")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		http.Error(w, "Ошибка при создании каталога для загрузки", http.StatusInternalServerError)
		return
	}

	file1Path := filepath.Join(uploadsDir, fmt.Sprintf("first_image_%d.jpg", time.Now().Unix()))
	file2Path := filepath.Join(uploadsDir, fmt.Sprintf("second_image_%d.jpg", time.Now().Unix()))

	saveFile := func(file multipart.File, path string) error {
		outFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		return err
	}

	if err := saveFile(file1, file1Path); err != nil {
		log.Printf("Ошибка при сохранении первого изображения: %v", err)
		http.Error(w, "Ошибка при сохранении первого изображения", http.StatusInternalServerError)
		return
	}

	if err := saveFile(file2, file2Path); err != nil {
		log.Printf("Ошибка при сохранении второго изображения: %v", err)
		http.Error(w, "Ошибка при сохранении второго изображения", http.StatusInternalServerError)
		return
	}

	log.Println(file1Path, file2Path)
	resultPath, err := facedetect.MatchFaces(file1Path, file2Path)
	if err != nil {
		log.Printf("Ошибка при обработке изображений: %v", err)
		http.Error(w, "Ошибка при обработке изображений", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Результат сопоставления сохранен в %s\n", resultPath)
}

func (h *Handler) PostFindMatchingFacesHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Фото сохранено на сервере: %s", inputPath)

	outputDir := "storage/detected_photo/"
	if err := facedetect.DetectFaces(inputPath, outputDir); err != nil {
		log.Printf("Ошибка при распознавании лиц: %v", err)
		http.Error(w, "Ошибка при распознавании лиц", http.StatusInternalServerError)
		return
	}

	log.Printf("Лица успешно распознаны на фото: %s", inputPath)

	faceDescriptors, err := facedetect.GetFaceDescriptors(outputDir)
	if err != nil {
		log.Printf("Ошибка при извлечении дескрипторов лиц: %v", err)
		http.Error(w, "Ошибка при извлечении дескрипторов лиц", http.StatusInternalServerError)
		return
	}

	log.Printf("Дескрипторы лиц извлечены успешно. Количество: %d", len(faceDescriptors))

	var allMatchingFaces []model.Face
	for _, descriptor := range faceDescriptors {
		matchingFaces, err := facedetect.FindMatchingFaces(descriptor, h.storage)
		if err != nil {
			log.Printf("Ошибка при поиске совпадений: %v", err)
			http.Error(w, fmt.Sprintf("Ошибка при поиске совпадений: %v", err), http.StatusInternalServerError)
			return
		}
		allMatchingFaces = append(allMatchingFaces, matchingFaces...)
	}
	payload := []struct {
		Path string `json:"path"`
	}{
		{Path: "https://i.imgur.com/WfbLlJ0.jpeg"},
		{Path: "https://i.imgur.com/WMbdvVW.jpeg"},
		{Path: "https://i.imgur.com/S9eTwSO.jpeg"},
	}
	json.NewEncoder(w).Encode(payload)
}
