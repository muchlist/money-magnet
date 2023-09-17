package main

import (
	"net/http"

	cyhand "github.com/muchlist/moneymagnet/business/category/handler"
	cyrepo "github.com/muchlist/moneymagnet/business/category/repo"
	cyserv "github.com/muchlist/moneymagnet/business/category/service"
	pthand "github.com/muchlist/moneymagnet/business/pocket/handler"
	ptrepo "github.com/muchlist/moneymagnet/business/pocket/repo"
	ptserv "github.com/muchlist/moneymagnet/business/pocket/service"
	reqhand "github.com/muchlist/moneymagnet/business/request/handler"
	reqrepo "github.com/muchlist/moneymagnet/business/request/repo"
	reqserv "github.com/muchlist/moneymagnet/business/request/service"
	spnhand "github.com/muchlist/moneymagnet/business/spend/handler"
	spnrepo "github.com/muchlist/moneymagnet/business/spend/repo"
	spnserv "github.com/muchlist/moneymagnet/business/spend/service"
	urhand "github.com/muchlist/moneymagnet/business/user/handler"
	urrepo "github.com/muchlist/moneymagnet/business/user/repo"
	urserv "github.com/muchlist/moneymagnet/business/user/service"
	"github.com/muchlist/moneymagnet/pkg/lrucache"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/muchlist/moneymagnet/pkg/mcrypto"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	// dependency
	jwt := mjwt.New(app.config.App.Secret)
	bcrypt := mcrypto.New()
	cache := lrucache.NewLRUCache()

	// middleware
	idempo := mid.NewIdempotencyMiddleware(cache)
	r.Use(mid.EndpoitnCounter)

	userRepo := urrepo.NewRepo(app.db, app.logger)
	pocketRepo := ptrepo.NewRepo(app.db, app.logger)
	categoryRepo := cyrepo.NewRepo(app.db, app.logger)
	requestRepo := reqrepo.NewRepo(app.db, app.logger)
	spendRepo := spnrepo.NewRepo(app.db, app.logger)

	userService := urserv.NewCore(app.logger, userRepo, bcrypt, jwt)
	userHandler := urhand.NewUserHandler(app.logger, app.validator, userService)

	pocketService := ptserv.NewCore(app.logger, pocketRepo, userRepo, categoryRepo)
	pocketHandler := pthand.NewPocketHandler(app.logger, app.validator, cache, pocketService)

	categoryService := cyserv.NewCore(app.logger, categoryRepo, pocketRepo)
	categoryHandler := cyhand.NewCatHandler(app.logger, app.validator, categoryService)

	requestService := reqserv.NewCore(app.logger, requestRepo, pocketRepo)
	requestHandler := reqhand.NewRequestHandler(app.logger, app.validator, requestService)

	spendService := spnserv.NewCore(app.logger, spendRepo, pocketRepo)
	spendHandler := spnhand.NewSpendHandler(app.logger, app.validator, cache, spendService)

	// swagger endpoint
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
	))

	// Endpoint with no auth required
	r.Get("/healthcheck", HealthCheckHandler)
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
			r.Get("/{id}", pocketHandler.GetByID)
			r.Get("/", pocketHandler.FindUserPocket)

			i := r.With(idempo.IdempotentCheck)
			i.Post("/", pocketHandler.CreatePocket)
			i.Patch("/{id}", pocketHandler.UpdatePocket)
		})

		r.Route("/categories", func(r chi.Router) {
			r.Post("/", categoryHandler.CreateCategory)
			r.Get("/from-pocket/{id}", categoryHandler.FindPocketCategory)
			r.Put("/{id}", categoryHandler.EditCategory)
			r.Delete("/{id}", categoryHandler.DeleteCategory)
		})

		r.Route("/request", func(r chi.Router) {
			r.Post("/{id}/action", requestHandler.ApproveOrRejectRequest)
			r.Post("/", requestHandler.CreateRequest)
			r.Get("/in", requestHandler.FindRequestByApprover)
			r.Get("/out", requestHandler.FindByRequester)
		})

		r.Route("/spends", func(r chi.Router) {
			r.Get("/from-pocket/{id}", spendHandler.FindSpend)
			r.Get("/{id}", spendHandler.GetByID)
			r.Post("/sync/{id}", spendHandler.SyncBalance)

			i := r.With(idempo.IdempotentCheck)
			i.Post("/", spendHandler.CreateSpend)
			i.Patch("/{id}", spendHandler.EditSpend)
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
