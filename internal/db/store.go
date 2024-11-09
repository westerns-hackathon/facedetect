package db

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"face-detection/internal/model"
	"fmt"
)

type Storage interface {
	AddFace(face model.Face) error
	GetAllFaces() ([]model.Face, error)
}

func (d *Database) AddFace(face model.Face) error {
	// Преобразуем metadata в JSON
	metadataJSON, err := json.Marshal(face.Metadata)
	if err != nil {
		return fmt.Errorf("не удалось преобразовать metadata в JSON: %v", err)
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(face.Descriptor)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать дескриптор: %v", err)
	}
	descriptorBytes := buf.Bytes()

	query := `INSERT INTO faces (metadata, descriptor, photo_path) VALUES ($1, $2, $3)`
	_, err = d.db.Exec(query, metadataJSON, descriptorBytes, face.PhotoPath)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении лица: %v", err)
	}
	return nil
}

func (d *Database) GetAllFaces() ([]model.Face, error) {
	query := `SELECT id, metadata, descriptor, photo_path FROM faces`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех лиц: %v", err)
	}
	defer rows.Close()

	var faces []model.Face
	for rows.Next() {
		var face model.Face
		var metadataBytes []byte
		var descriptorBytes []byte
		if err := rows.Scan(&face.ID, &metadataBytes, &descriptorBytes, &face.PhotoPath); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %v", err)
		}

		if err := json.Unmarshal(metadataBytes, &face.Metadata); err != nil {
			return nil, fmt.Errorf("не удалось декодировать metadata: %v", err)
		}

		if err := deserializeDescriptor(descriptorBytes, &face.Descriptor); err != nil {
			return nil, fmt.Errorf("не удалось десериализовать дескриптор: %v", err)
		}

		faces = append(faces, face)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %v", err)
	}

	return faces, nil
}

func deserializeDescriptor(data []byte, descriptor *[]float32) error {

	decoder := gob.NewDecoder(bytes.NewReader(data))

	if err := decoder.Decode(descriptor); err != nil {
		return fmt.Errorf("не удалось десериализовать дескриптор: %v", err)
	}

	return nil
}
