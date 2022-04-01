package main

import (
	"expvar"
	"net/http"
	"runtime"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/app/api/handler"
	"github.com/muchlist/moneymagnet/bussines/core/user/userrepo"
	"github.com/muchlist/moneymagnet/bussines/core/user/userservice"
	"github.com/muchlist/moneymagnet/bussines/sys/mjwt"
	"github.com/muchlist/moneymagnet/foundation/mcrypto"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	// dependency
	jwt := mjwt.New(app.config.secret)
	bcrypt := mcrypto.New()

	userRepo := userrepo.NewRepo(app.db)
	userService := userservice.NewService(app.logger, userRepo, bcrypt, jwt)
	userHandler := handler.NewUserHandler(app.logger, userService)

	router.Get("/healthcheck", handler.HealthCheckHandler)
	router.Post("/login", userHandler.Login)
	router.Post("/register", userHandler.Register)

	// setup exvar for monitoring
	setupExpvar(app.db)
	router.Mount("/debug/vars", expvar.Handler())

	return router
}

// =============================================================================

// setupExpvar setup exvar for monitoring
func setupExpvar(db *pgxpool.Pool) {
	expvar.NewString("api_version").Set(version)
	expvar.Publish("api_timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("database", expvar.Func(func() interface{} {
		stat := db.Stat()
		return map[string]interface{}{
			"conn_max":       stat.MaxConns(),
			"conn_idle":      stat.IdleConns(),
			"conn_in_use":    stat.TotalConns(),
			"acquire_total":  stat.AcquireCount(),
			"acquire_cancel": stat.CanceledAcquireCount(),
		}
	}))
}
