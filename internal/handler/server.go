package handler

import (
	"face-detection/internal/config"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	h   *Handler
	cfg *config.Config
}

func NewServer(h *Handler) *Server {
	return &Server{h: h}
}

func (s *Server) Run(host, port string) error {
	mux := http.NewServeMux()
	log.Printf("server in running on %s:%s", host, port)
	mux.HandleFunc("POST /v1/app/photo", s.h.PostPhotoHandler)
	mux.HandleFunc("POST /v1/app/photo/match", s.h.PostFaceMatchHandler)
	mux.HandleFunc("POST /v1/app/photo/find", s.h.PostFindMatchingFacesHandler)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), mux); err != nil {
		return err
	}
	return nil
}
