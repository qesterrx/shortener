package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qesterrx/shortener/internal/config"
	"github.com/qesterrx/shortener/internal/handler"
	"github.com/qesterrx/shortener/internal/logger"
	"github.com/qesterrx/shortener/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	logger.InitLogger()
	config := config.ParseParams()

	fmt.Println("ServerHost=", config.ServerHost.String())
	fmt.Println("ServerRedirect=", config.ServerRedirect.String())
	fmt.Println("FileStorage=", config.FileStorage)

	pool, err := pgxpool.New(context.TODO(), config.DatabaseDSN)
	if err != nil {
		logger.Log.Error().Msg("Error connect to DB")
		panic(err)
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	storage, err := service.NewURLShortnerStorage(config.FileStorage)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	return http.ListenAndServe(config.ServerHost.String(), handler.Router(storage, config, db))
}
