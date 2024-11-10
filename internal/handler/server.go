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
	directory := "storage/"

	fs := http.FileServer(http.Dir(directory))
	mux.Handle("/dir/", http.StripPrefix("/dir/", fs))
	mux.HandleFunc("/v1/app/photo", s.h.PostPhotoHandler)
	mux.HandleFunc("/v1/app/photo/match", s.h.PostFaceMatchHandler)
	mux.HandleFunc("/v1/app/photo/find", s.h.PostFindMatchingFacesHandler)

	handler := enableCORS(mux)

	log.Printf("server is running on %s:%s", host, port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), handler); err != nil {
		return err
	}
	return nil
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
