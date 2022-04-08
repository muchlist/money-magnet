package main

import (
	"net/http"

	"github.com/muchlist/moneymagnet/app/api/handler"
	"github.com/muchlist/moneymagnet/bussines/core/user/userrepo"
	"github.com/muchlist/moneymagnet/bussines/core/user/userservice"
	"github.com/muchlist/moneymagnet/bussines/sys/mid"
	"github.com/muchlist/moneymagnet/bussines/sys/mjwt"
	"github.com/muchlist/moneymagnet/foundation/mcrypto"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	// dependency
	jwt := mjwt.New(app.config.secret)
	bcrypt := mcrypto.New()

	userRepo := userrepo.NewRepo(app.db)
	userService := userservice.NewService(app.logger, userRepo, bcrypt, jwt)
	userHandler := handler.NewUserHandler(app.logger, userService)

	// Endpoint with no auth required
	r.Get("/healthcheck", handler.HealthCheckHandler)
	r.Post("/login", userHandler.Login)

	// Endpoint with fresh auth admin
	r.Group(func(r chi.Router) {
		r.Use(mid.RequiredFreshRoles("admin"))
		r.Post("/register", userHandler.Register)
	})

	// Endpoint with auth
	r.Group(func(r chi.Router) {
		r.Use(mid.RequiredRoles())
		r.Get("/profile", userHandler.Profile)
	})

	return r
}

// =============================================================================
