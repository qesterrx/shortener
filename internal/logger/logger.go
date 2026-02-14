package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

type respData struct {
	statusCode int
	size       int
}

type loggedResponseWriter struct {
	http.ResponseWriter
	data respData
}

func (lrw *loggedResponseWriter) Write(msg []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(msg)
	lrw.data.size = size
	return size, err
}

func (lrw *loggedResponseWriter) WriteHeader(statusCode int) {
	lrw.data.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

// Singleton
func InitLogger() {

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	Log = log
}

// Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
// Сведения об ответах должны содержать код статуса и размер содержимого ответа.
func HandlerLogger(h http.HandlerFunc) http.HandlerFunc {
	loggedFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := loggedResponseWriter{ResponseWriter: w, data: respData{}}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.Info().
			Str("uri", r.RequestURI).
			Str("method", r.Method).
			Str("duration", duration.String()).
			Int("status", lw.data.statusCode).
			Int("size", lw.data.size).
			Send()
	}

	return http.HandlerFunc(loggedFunc)
}
