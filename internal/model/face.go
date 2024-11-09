package model

type Face struct {
	ID         int       `json:"id"`
	Metadata   []string  `json:"metadata"`
	Descriptor []float32 `json:"descriptor"`
	PhotoPath  string    `json:"photo_path"`
	Distance   float32   `json:"distance,omitempty"` // добавлено для хранения расстояния при поиске
}
