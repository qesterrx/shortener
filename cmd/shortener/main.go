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

	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handler.ShortURL(storage, config))
	mux.HandleFunc(`/{id}`, handler.GetFullURL(storage, config))

	return http.ListenAndServe(config.HOST, mux)
}
