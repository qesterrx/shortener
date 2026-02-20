package main

import (
	"fmt"
	"net/http"

	"github.com/qesterrx/shortener/internal/config"
	"github.com/qesterrx/shortener/internal/handler"
	"github.com/qesterrx/shortener/internal/logger"
	"github.com/qesterrx/shortener/internal/service"
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

	storage, err := service.NewURLShortnerStorage(config.FileStorage)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	return http.ListenAndServe(config.ServerHost.String(), handler.Router(storage, config))
}
