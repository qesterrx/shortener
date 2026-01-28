package main

import (
	"net/http"

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

	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handler.ShortURL(storage))
	mux.HandleFunc(`/{id}`, handler.GetFullURL(storage))

	return http.ListenAndServe(`:8080`, mux)
}
