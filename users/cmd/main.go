package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/madrabit/mini-market/users/internal"
	"github.com/madrabit/mini-market/users/internal/common"
	"github.com/madrabit/mini-market/users/internal/database"
	"github.com/madrabit/mini-market/users/internal/validator"
	"github.com/madrabit/mini-market/users/internal/web"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := common.Load()
	if err != nil {
		fmt.Println("config load error, %w", err)
		//падение конфига как то отлогировать
	}
	logger := common.NewLogger(cfg)
	vld := validator.New()
	db := database.ConnectDbWithCfg(cfg)
	repository := internal.NewRepository(db)
	service := internal.NewService(repository, vld)
	controllerUsers := internal.NewControllerUsers(service, logger)
	controllerRoles := internal.NewControllerRoles(service, logger)
	server := web.NewServer()
	server.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Mount("/users", controllerUsers.Routes())
			r.Mount("/roles", controllerRoles.Routes())
		})
	})
	err = http.ListenAndServe(":8080", server.Router)
	if err != nil {
		return
	}
}
