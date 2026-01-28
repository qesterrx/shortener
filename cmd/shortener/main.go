package main

import (
	"net/http"

	"github.com/qesterrx/shortener/internal/config"
	"github.com/qesterrx/shortener/internal/handler"
	"github.com/qesterrx/shortener/internal/service"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	storage := &service.URLShortnerStorage{Storage: make(map[string]string)}
	config := &config.Config{HOST: "localhost:8080"}

	return http.ListenAndServe(config.HOST, handler.Router(storage, config))
}
