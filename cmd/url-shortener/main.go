package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"url-shortener/internal/config"
	"url-shortener/pkg/logger/sl"
	"url-shortener/storage/sqlite"

	"url-shortener/internal/http-server/router"
)

func main() {
	cfg := config.MustLoad()

	// Устанавливаем дефолтные параметры логера в зависимости от уровня окружения
	config.SetupLogger(cfg.EnvLevel, os.Stdout)

	log := slog.Default().With(slog.String("env", cfg.EnvLevel))

	log.Info("Init server", slog.String("address", cfg.Address))
	log.Debug("Debug mode ON")

	storage, err := sqlite.New(context.Background(), cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

	router := router.GetRouter(router.NewRouterApp(storage, cfg))

	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		slog.Error("Ошибка запуска сервера", "err", err.Error())
	}
}
