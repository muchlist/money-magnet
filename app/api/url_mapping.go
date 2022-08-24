package main

import (
	"net/http"

	cyrepo "github.com/muchlist/moneymagnet/business/category/repo"
	cyserv "github.com/muchlist/moneymagnet/business/category/service"
	ptrepo "github.com/muchlist/moneymagnet/business/pocket/repo"
	ptserv "github.com/muchlist/moneymagnet/business/pocket/service"
	urrepo "github.com/muchlist/moneymagnet/business/user/repo"
	urserv "github.com/muchlist/moneymagnet/business/user/service"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/mjwt"

	"github.com/muchlist/moneymagnet/app/api/handler"
	"github.com/muchlist/moneymagnet/pkg/mcrypto"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	// dependency
	jwt := mjwt.New(app.config.secret)
	bcrypt := mcrypto.New()

	userRepo := urrepo.NewRepo(app.db)
	userService := urserv.NewCore(app.logger, userRepo, bcrypt, jwt)
	userHandler := handler.NewUserHandler(app.logger, userService)

	pocketRepo := ptrepo.NewRepo(app.db)
	pocketService := ptserv.NewCore(app.logger, pocketRepo, userRepo)
	pocketHandler := handler.NewPocketHandler(app.logger, pocketService)

	categoryRepo := cyrepo.NewRepo(app.db)
	categoryService := cyserv.NewCore(app.logger, categoryRepo)
	categoryHandler := handler.NewCatHandler(app.logger, categoryService)

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
		r.Route("/user", func(r chi.Router) {
			r.Get("/profile", userHandler.Profile)
			r.Get("/{strID}", userHandler.GetByID)
			r.Get("/", userHandler.FindByName)
			r.Post("/fcm/{strID}", userHandler.UpdateFCM)
		})

		r.Route("/pockets", func(r chi.Router) {
			r.Post("/", pocketHandler.CreatePocket)
			r.Get("/{id}", pocketHandler.GetByID)
			r.Get("/", pocketHandler.FindUserPocket)
			r.Put("/", pocketHandler.RenamePocket)
		})

		r.Route("/categories", func(r chi.Router) {
			r.Post("/", categoryHandler.CreateCategory)
			r.Get("/from-pocket/{id}", categoryHandler.FindPocketCategory)
			r.Put("/", categoryHandler.EditCategory)
			r.Delete("/{strID}", categoryHandler.DeleteCategory)
		})
	})

	// Endpoint with fresh auth
	r.Group(func(r chi.Router) {
		r.Use(mid.RequiredFreshRoles())
		r.Patch("/user/profile", userHandler.EditSelfUser)
	})

	return r
}

// =============================================================================
