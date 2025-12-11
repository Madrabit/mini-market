package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/madrabit/mini-market/users/internal"
	"github.com/madrabit/mini-market/users/internal/common"
	"github.com/madrabit/mini-market/users/internal/database"
	"github.com/madrabit/mini-market/users/internal/validator"
	"github.com/madrabit/mini-market/users/internal/web"
	gracefulshutdown "github.com/quii/go-graceful-shutdown"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := common.Load()
	if err != nil {
		log.Fatal("config load error, %w", err)
	}
	logger := common.NewLogger(cfg)
	db := database.ConnectDbWithCfg(cfg)
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Error("failed to close db")
		}
	}()
	server := build(db, logger)
	httpServer := &http.Server{Addr: cfg.Server.Port, Handler: server.Router}
	ctx := context.Background()
	srv := gracefulshutdown.NewServer(httpServer)
	if err := srv.ListenAndServe(ctx); err != nil {
		logger.Fatal("didnt shutdown gracefully, some responses may have been lost")
	}
	logger.Info("shutdown gracefully! all responses were sent")
}

func build(db *sqlx.DB, logger *common.Logger) *web.Server {
	server := web.NewServer()
	vld := validator.New()
	repository := internal.NewRepository(db)
	roleService := internal.NewRoleService(repository, vld)
	userService := internal.NewUserService(repository, roleService, vld)
	controllerUsers := internal.NewControllerUsers(userService, logger)
	controllerRoles := internal.NewControllerRoles(roleService, logger)
	server.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Mount("/users", controllerUsers.Routes())
			r.Mount("/roles", controllerRoles.Routes())
		})
	})
	return server
}
