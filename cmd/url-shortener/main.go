package main

import (
	"fmt"
	"http-app/internal/config"
	"http-app/internal/lib/logger/sl"
	"http-app/internal/storage/sqlite"
	"log/slog"
	"os"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	//TODO: init config: cleanenv
	cfg := config.MustLoad()

	fmt.Println(cfg)

	//TODO: init logger: slog

	log := setupLogger(cfg.Env)

	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	//TODO: init storage: sqlite

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	//TODO: init router: chi, render

	router := chi.NewRouter()

	//middleware

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	//router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}
	return log

}
