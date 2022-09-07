package main

import (
	"net/http"

	cyrepo "github.com/muchlist/moneymagnet/business/category/repo"
	cyserv "github.com/muchlist/moneymagnet/business/category/service"
	ptrepo "github.com/muchlist/moneymagnet/business/pocket/repo"
	ptserv "github.com/muchlist/moneymagnet/business/pocket/service"
	reqrepo "github.com/muchlist/moneymagnet/business/request/repo"
	reqserv "github.com/muchlist/moneymagnet/business/request/service"
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
	userHandler := handler.NewUserHandler(app.logger, app.validator, userService)

	pocketRepo := ptrepo.NewRepo(app.db)
	pocketService := ptserv.NewCore(app.logger, pocketRepo, userRepo)
	pocketHandler := handler.NewPocketHandler(app.logger, app.validator, pocketService)

	categoryRepo := cyrepo.NewRepo(app.db)
	categoryService := cyserv.NewCore(app.logger, categoryRepo)
	categoryHandler := handler.NewCatHandler(app.logger, app.validator, categoryService)

	requestRepo := reqrepo.NewRepo(app.db)
	requestService := reqserv.NewCore(app.logger, requestRepo, pocketRepo)
	requestHandler := handler.NewRequestHandler(app.logger, app.validator, requestService)

	// Endpoint with no auth required
	r.Get("/healthcheck", handler.HealthCheckHandler)
	r.Post("/user/login", userHandler.Login)
	r.Post("/user/refresh", userHandler.RefreshToken)

	// Endpoint with fresh auth admin
	r.Group(func(r chi.Router) {
		r.Use(mid.RequiredFreshRoles("admin"))
		r.Post("/register", userHandler.Register)
		r.Patch("/edit-user/{id}", userHandler.EditUser)
		r.Delete("/user/{id}", userHandler.DeleteUser)
	})

	// Endpoint with auth
	r.Group(func(r chi.Router) {
		r.Use(mid.RequiredRoles())
		r.Route("/user", func(r chi.Router) {
			r.Get("/profile", userHandler.Profile)
			r.Get("/{id}", userHandler.GetByID)
			r.Get("/", userHandler.FindByName)
			r.Post("/fcm/{id}", userHandler.UpdateFCM)
		})

		r.Route("/pockets", func(r chi.Router) {
			r.Post("/", pocketHandler.CreatePocket)
			r.Get("/{id}", pocketHandler.GetByID)
			r.Get("/", pocketHandler.FindUserPocket)
			r.Put("/rename", pocketHandler.RenamePocket)
		})

		r.Route("/categories", func(r chi.Router) {
			r.Post("/", categoryHandler.CreateCategory)
			r.Get("/from-pocket/{id}", categoryHandler.FindPocketCategory)
			r.Put("/", categoryHandler.EditCategory)
			r.Delete("/{id}", categoryHandler.DeleteCategory)
		})

		r.Route("/request", func(r chi.Router) {
			r.Post("/{id}/action", requestHandler.ApproveOrRejectRequest)
			r.Post("/", requestHandler.CreateRequest)
			r.Get("/in", requestHandler.FindRequestByApprover)
			r.Get("/out", requestHandler.FindByRequester)
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
