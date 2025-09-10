package web

import "github.com/go-chi/chi/v5"

type Server struct {
	Router chi.Router
}

func NewServer() *Server {
	r := chi.NewRouter()
	return &Server{Router: r}
}
