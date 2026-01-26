package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/qesterrx/shortener/internal/service"
)

func Shortner(storage *service.URLShortnerStorage) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {

			path := strings.Trim(r.URL.Path, `/`)
			parts := strings.Split(path, `/`)

			if len(parts) != 1 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fullUrl := storage.Get(parts[0])

			w.Header().Set("Location", fullUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)

			return

		}

		if r.Method == http.MethodPost {

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

			shortUrl := "http://localhost:8080/" + storage.Set(string(body))

			storage.Show()
			//Remember
			w.Header().Set("Content-Type", "text/plain")

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortUrl))

			return

		}

		w.WriteHeader(http.StatusBadRequest)
		return
	})

}
