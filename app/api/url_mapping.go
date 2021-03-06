package main

import (
	"net/http"

	"github.com/muchlist/moneymagnet/app/api/handler"
	"github.com/muchlist/moneymagnet/bussines/core/pocket/ptrepo"
	"github.com/muchlist/moneymagnet/bussines/core/pocket/ptservice"
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

	pocketRepo := ptrepo.NewRepo(app.db)
	pocketService := ptservice.NewService(app.logger, pocketRepo, userRepo)
	pocketHandler := handler.NewPocketHandler(app.logger, pocketService)

	// Endpoint with no auth required
	r.Get("/healthcheck", handler.HealthCheckHandler)
	r.Post("/user/login", userHandler.Login)
	r.Post("/user/refresh", userHandler.RefreshToken)

	// Endpoint with fresh auth admin
	r.Group(func(r chi.Router) {
		r.Use(mid.RequiredFreshRoles("admin"))
		r.Post("/register", userHandler.Register)
		r.Patch("/edit-user/{strID}", userHandler.EditUser)
		r.Delete("/user/{strID}", userHandler.DeleteUser)
	})

	// Endpoint with auth
	r.Group(func(r chi.Router) {
		r.Use(mid.RequiredRoles())
		r.Get("/user/profile", userHandler.Profile)
		r.Get("/user/{strID}", userHandler.GetByID)
		r.Get("/user", userHandler.FindByName)
		r.Post("/user/fcm/{strID}", userHandler.UpdateFCM)

		r.Post("/pockets", pocketHandler.CreatePocket)
		r.Get("/pockets/{id}", pocketHandler.GetByID)
		r.Get("/pockets", pocketHandler.FindUserPocket)
		r.Put("/rename-pocket", pocketHandler.RenamePocket)
	})

	// Endpoint with fresh auth
	r.Group(func(r chi.Router) {
		r.Use(mid.RequiredFreshRoles())
		r.Patch("/user/profile", userHandler.EditSelfUser)
	})

	return r
}

// =============================================================================
