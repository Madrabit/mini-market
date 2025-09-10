package web

import "github.com/go-chi/chi/v5"

type Server struct {
	Router chi.Router
}

func NewServer() *Server {
	router := chi.NewRouter()
	return &Server{Router: router}
}
