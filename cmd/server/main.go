package main

import (
	"context"
	"github.com/51mans0n/effective-mobile-task/internal/config"
	"github.com/51mans0n/effective-mobile-task/internal/http/handler"
	"github.com/51mans0n/effective-mobile-task/internal/migration"
	"github.com/51mans0n/effective-mobile-task/internal/repository"
	"github.com/51mans0n/effective-mobile-task/internal/service"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mymw "github.com/51mans0n/effective-mobile-task/internal/http/middleware"
	mylog "github.com/51mans0n/effective-mobile-task/internal/logger"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	db, err := sqlx.Open("pgx", cfg.DBDSN)
	if err != nil {
		logger.Fatal("db connect", zap.Error(err))
	}
	if err := db.Ping(); err != nil {
		logger.Fatal("db ping", zap.Error(err))
	}

	// auto-migrate
	if os.Getenv("AUTO_MIGRATE") == "true" {
		if err := migration.Up(db); err != nil {
			logger.Fatal("migrate", zap.Error(err))
		}
	}
	repo := repository.NewPeopleRepo(db)
	_ = repo
	svc := service.New(repo)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)   // 500 handler
	r.Use(mymw.ZapLogger(logger)) // custom zap logger

	// health-check
	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	r.Post("/people", handler.Create(svc))
	r.Get("/people", handler.List(svc))

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
