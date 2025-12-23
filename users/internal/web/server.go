package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

type Server struct {
	Router chi.Router
}

func NewServer(reg prometheus.Registerer) *Server {
	r := chi.NewRouter()
	Init(reg)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(Prometheus)
	return &Server{Router: r}
}
