package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"go.uber.org/zap"
)

type webServer struct {
	logger      mlogger.Logger
	wg          sync.WaitGroup // waitgroup for gracefully shutdown background process
	port        int
	env         string
	serviceName string
}

func New(logger mlogger.Logger, port int, env string, serviceName string) *webServer {
	return &webServer{
		logger:      logger,
		wg:          sync.WaitGroup{},
		port:        port,
		env:         env,
		serviceName: serviceName,
	}
}

func (ws *webServer) Serve(injectRoute http.Handler) error {
	// Declare a HTTP server using the same settings as in our main() function.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", ws.port),
		Handler:      ws.setupRoutes(injectRoute),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// to receive any errows returned
	// by the graceful Shutdown() function.
	shutdownError := make(chan error)

	// gracefully shutdown.
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ws.logger.Info("shutting down server", zap.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		ws.logger.Info("completing background tasks")

		ws.wg.Wait()
		shutdownError <- nil
	}()

	ws.logger.Info("starting server", zap.String("addr", srv.Addr), zap.String("env", ws.env))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	ws.logger.Info("stopped server", zap.String("addr", srv.Addr))
	return nil
}
