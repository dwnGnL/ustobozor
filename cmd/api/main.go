package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/barzurustami/bozor/internal/app"
	"github.com/barzurustami/bozor/internal/config"
	"github.com/barzurustami/bozor/internal/db"
	"github.com/barzurustami/bozor/internal/graphql"
	"github.com/barzurustami/bozor/internal/logger"
	"github.com/barzurustami/bozor/internal/middleware"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.App.LogLevel)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	pool, err := db.Connect(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatal("db connect failed", zap.Error(err))
	}
	defer pool.Close()

	repos := app.NewRepositories(pool)
	services := app.NewServices(cfg, repos, log)
	resolver := app.NewResolver(services, repos)

	gqlServer := graphql.NewServer(resolver, cfg.Upload.MaxSizeBytes)

	mux := http.NewServeMux()
	mux.Handle("/graphql", graphql.MaxBytes(cfg.Upload.MaxSizeBytes, gqlServer))
	mux.Handle("/", playground.Handler("Bozor GraphQL", "/graphql"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.Upload.Dir))))

	h := middleware.RequestID(mux)
	h = middleware.Logging(log)(h)
	h = middleware.Auth(services.JWT)(h)

	srv := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Info("server started", zap.String("port", cfg.App.Port))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown error", zap.Error(err))
	}
}
