package main

import (
	"context"
	"github.com/51mans0n/effective-mobile-task/internal/config"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mymw "github.com/51mans0n/effective-mobile-task/internal/http/middleware"
	mylog "github.com/51mans0n/effective-mobile-task/internal/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger, err := mylog.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flush

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)   // 500 handler
	r.Use(mymw.ZapLogger(logger)) // custom zap logger

	// health-check
	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// graceful-shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	logger.Info("server start", zap.String("addr", cfg.AppPort))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("listen", zap.Error(err))
	}
}
