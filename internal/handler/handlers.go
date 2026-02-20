package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/qesterrx/shortener/internal/compression"
	"github.com/qesterrx/shortener/internal/config"
	"github.com/qesterrx/shortener/internal/logger"
	"github.com/qesterrx/shortener/internal/model"
	"github.com/qesterrx/shortener/internal/service"
)

func ShortURL(storage *service.URLShortnerStorage, config *config.Configuration) http.HandlerFunc {

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
		w.Write([]byte(fmt.Sprintf("http://%s/%s", config.ServerRedirect.String(), shortURL)))

	})

}

func GetFullURL(storage *service.URLShortnerStorage, config *config.Configuration) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortURL := r.PathValue("id")
		//fmt.Printf("shortURL=%s\n", shortURL)

		if shortURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fullURL, err := storage.Get(shortURL)
		//fmt.Printf("fullURL=%s\n", fullURL)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	})

}

func ShortJSON(storage *service.URLShortnerStorage, config *config.Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost ||
			r.Header.Get("Content-Type") != "application/json" ||
			r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var req model.ShortenURLReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		short, err := storage.Set(string(req.URL))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res := model.ShortenURLRes{
			Goto: fmt.Sprintf("http://%s/%s", config.ServerRedirect.String(), short),
		}

		body, err := json.Marshal(&res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		w.Write(body)
	}
}

func Router(storage *service.URLShortnerStorage, config *config.Configuration) chi.Router {
	r := chi.NewRouter()

	r.Post(`/api/shorten`, compression.HandlerGzipCompress(logger.HandlerLogger(ShortJSON(storage, config))))
	r.Post(`/`, compression.HandlerGzipCompress(logger.HandlerLogger(ShortURL(storage, config))))
	r.Get(`/{id}`, compression.HandlerGzipCompress(logger.HandlerLogger(GetFullURL(storage, config))))

	return r
}
