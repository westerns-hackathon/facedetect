package model

type Face struct {
	ID       int      `json:"face_id"`
	Metadata []string `json:"metadata"`
}
