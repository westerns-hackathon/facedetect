package handler

import (
	"log"
	"net/http"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {
	mux := http.NewServeMux()
	log.Printf("server in running")
	http.ListenAndServe(":8080", mux)
}
