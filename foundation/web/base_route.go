package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (ws *webServer) setupRoutes(injectRoute http.Handler) http.Handler {
	router := chi.NewRouter()

	// convert notFoundResponse to http handler and set it as the custom error handler for 404 notfound
	router.NotFound(NotFoundResponse)
	// convert methodNotAllowedResponse to http handler and set it as the custom error handler for 405 method not allowed
	router.MethodNotAllowed(MethodNotAllowedResponse)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(midLogger(ws.logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))

	router.Mount("/", injectRoute)

	return router
}
