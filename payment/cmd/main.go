package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/madrabit/mini-market/payment/internal/web"
	"net/http"
)

func main() {
	controller := internal.NewController()
	server := web.NewServer()
	server.Router.Route("/api", func(r chi.Router) {
		r.Route("v1", func(r chi.Router) {
			r.Mount("/payments", controller.Routes())
		})
	})
	err := http.ListenAndServe("8084", server.Router)
	if err != nil {
		return
	}
}
