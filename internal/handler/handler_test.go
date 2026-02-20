package handler

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/qesterrx/shortener/internal/config"
	"github.com/qesterrx/shortener/internal/model"
	"github.com/qesterrx/shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortURL(t *testing.T) {
	config := &config.Configuration{
		ServerHost:     config.NetAddress{Host: "localhost", Port: 8080},
		ServerRedirect: config.NetAddress{Host: "localhost", Port: 8080},
	}

	storage := &service.URLShortnerStorage{Storage: make(map[string]*model.ShortenURL)}
	duplicateString := "hash for this string was added to storage"
	storage.Storage[service.GetHash(duplicateString)] = &model.ShortenURL{UUID: 1, ShortURL: service.GetHash(duplicateString), OriginalURL: "duplicate string"}

	//storage.Show()

	type want struct {
		statusCode  int
		contentType string
		body        string
	}

	tests := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			name:   "Correct empty body",
			method: http.MethodPost,
			body:   "",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				body:        fmt.Sprintf("http://%s/%s", config.ServerRedirect.String(), service.GetHash("")),
			},
		},
		{
			name:   "Correct not empty body",
			method: http.MethodPost,
			body:   "https://habr.com/ru/articles/550352/",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				body:        fmt.Sprintf("http://%s/%s", config.ServerRedirect.String(), service.GetHash("https://habr.com/ru/articles/550352/")),
			},
		},
		{
			name:   "Error duplicate hash",
			method: http.MethodPost,
			body:   duplicateString,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "Error wrong method GET",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "Error wrong method PUT",
			method: http.MethodPut,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	handler := ShortURL(storage, config)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()

			handler(w, r)

			assert.Equal(t, test.want.statusCode, w.Code, fmt.Sprintf("Code in response different: want %d got %d", test.want.statusCode, w.Code))
			if test.want.body != "" {
				assert.Equal(t, test.want.contentType, w.Header().Get("Content-Type"), fmt.Sprintf("Content-Type in response different: want %s got %s", test.want.contentType, w.Header().Get("Content-Type")))
				body := w.Body.String()
				assert.Equal(t, test.want.body, body, fmt.Sprintf("Body in response different: want %s got %s", test.want.body, body))
			}
		})
	}

}

func TestShortJSON(t *testing.T) {
	config := &config.Configuration{
		ServerHost:     config.NetAddress{Host: "localhost", Port: 8080},
		ServerRedirect: config.NetAddress{Host: "localhost", Port: 8080},
	}

	storage := &service.URLShortnerStorage{Storage: make(map[string]*model.ShortenURL)}
	req := model.ShortenURLReq{URL: "https://habr.com/ru/articles/550352/"}
	reqJSON, _ := json.Marshal(&req)
	res := model.ShortenURLRes{Goto: fmt.Sprintf("http://%s/%s", config.ServerRedirect.String(), service.GetHash("https://habr.com/ru/articles/550352/"))}
	resJSON, _ := json.Marshal(&res)

	type want struct {
		statusCode  int
		contentType string
		body        string
	}

	tests := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			name:   "empty body",
			method: http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "NO JSON body",
			method: http.MethodPost,
			body:   "hello",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "Correct not empty body",
			method: http.MethodPost,
			body:   string(reqJSON),
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
				body:        string(resJSON),
			},
		},
		{
			name:   "Error wrong method GET",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "Error wrong method PUT",
			method: http.MethodPut,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "Error wrong method DELETE",
			method: http.MethodDelete,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	handler := ShortJSON(storage, config)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "/api/shorten", strings.NewReader(test.body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler(w, r)

			assert.Equal(t, test.want.statusCode, w.Code, fmt.Sprintf("Code in response different: want %d got %d", test.want.statusCode, w.Code))
			if test.want.body != "" {
				assert.Equal(t, test.want.contentType, w.Header().Get("Content-Type"), fmt.Sprintf("Content-Type in response different: want %s got %s", test.want.contentType, w.Header().Get("Content-Type")))
				body := w.Body.String()
				assert.Equal(t, test.want.body, body, fmt.Sprintf("Body in response different: want %s got %s", test.want.body, body))
			}
		})
	}

}

func TestGetFullURL(t *testing.T) {

	config := &config.Configuration{
		ServerHost:     config.NetAddress{Host: "localhost", Port: 8080},
		ServerRedirect: config.NetAddress{Host: "localhost", Port: 8080},
	}

	storage := &service.URLShortnerStorage{Storage: make(map[string]*model.ShortenURL)}
	savedValue := "test"
	storage.Storage[service.GetHash(savedValue)] = &model.ShortenURL{UUID: 1, ShortURL: service.GetHash(savedValue), OriginalURL: savedValue}

	//storage.Show()

	mux := http.NewServeMux()
	mux.HandleFunc("/", GetFullURL(storage, config))
	mux.HandleFunc("/{id}", GetFullURL(storage, config))

	type want struct {
		statusCode int
		location   string
	}

	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{
			name:   "Correct saved value",
			url:    fmt.Sprintf("/%s", service.GetHash(savedValue)),
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   savedValue,
			},
		},
		{
			name:   "Error wrong method POST",
			method: http.MethodPost,
			url:    "/",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//fmt.Println(test.url)
			r := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, r)

			require.Equal(t, test.want.statusCode, w.Code, fmt.Sprintf("Code in response different: want %d got %d", test.want.statusCode, w.Code))
			if test.want.statusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, test.want.location, w.Header().Get("Location"), fmt.Sprintf("Location in response different: want %s got %s", test.want.location, w.Header().Get("Location")))
			}
		})
	}

}
