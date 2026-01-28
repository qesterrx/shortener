package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/qesterrx/shortener/internal/config"
	"github.com/qesterrx/shortener/internal/service"
)

func ShortURL(storage *service.URLShortnerStorage, config *config.Config) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		body, err := io.ReadAll(r.Body) //body is []byte

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//Remember new value
		shortURL, err := storage.Set(string(body))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//storage.Show() //---Отладка

		w.Header().Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%s%s", config.HOST, shortURL)))

	})

}

func GetFullURL(storage *service.URLShortnerStorage, config *config.Config) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortURL := r.PathValue("id")

		if shortURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fullURL, err := storage.Get(shortURL)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	})

}
