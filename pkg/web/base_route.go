package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
)

func (ws *webServer) setupRoutes(injectRoute http.Handler) http.Handler {
	router := chi.NewRouter()

	// convert notFoundResponse to http handler and set it as the custom error handler for 404 notfound
	router.NotFound(NotFoundResponse)
	// convert methodNotAllowedResponse to http handler and set it as the custom error handler for 405 method not allowed
	router.MethodNotAllowed(MethodNotAllowedResponse)

	router.Use(otelchi.Middleware(ws.serviceName, otelchi.WithChiRoutes(router)))
	router.Use(requestID)
	router.Use(middleware.RealIP)
	router.Use(midLogger(ws.logger))
	router.Use(panicRecovery(ws.logger))
	// router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))

	router.Mount("/", injectRoute)

	return router
}
