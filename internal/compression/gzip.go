package compression

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

//---------------------------------------------------------------------compressWriter

// в данном случае мы используем не встраивание а композицию, поэтому надо определить все методы интерфейса под который мы мимикрируем
type compressWriter struct {
	w  http.ResponseWriter //этот интерфейс у нас есть напрямую в хенделее
	zw *gzip.Writer
}

// Дополнительные процедуры типа конструктор/деструктор
func newCompressWriter(dst http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  dst,
		zw: gzip.NewWriter(dst),
	}
}

func (cw *compressWriter) Close() error {
	return cw.zw.Close()
}

// Реализация интерфейса
func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	return cw.zw.Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
	cw.w.WriteHeader(statusCode)
}

//---------------------------------------------------------------------compressReader

// в данном случае мы используем не встраивание а композицию, поэтому надо определить все методы интерфейса под который мы мимикрируем
type compressReader struct {
	r  io.ReadCloser //этот интерфейс реализует *http.request.body
	zr *gzip.Reader
}

//Дополнительные процедуры типа конструктор

func NewCompressReader(src io.ReadCloser) (*compressReader, error) {

	nr, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  src,
		zr: nr,
	}, nil
}

// Реализация интерфейса
func (cr *compressReader) Read(p []byte) (n int, err error) {
	return cr.zr.Read(p)
}

func (cr *compressReader) Close() error {
	err := cr.r.Close()
	if err != nil {
		return err
	}

	return cr.zr.Close()
}

//---------------------------------------------------------------------Middleware

func HandlerGzipCompress(h http.HandlerFunc) http.HandlerFunc {
	funcCompress := func(w http.ResponseWriter, r *http.Request) {

		finw := w

		//Подменяем райтер

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			cw := newCompressWriter(w)
			finw = cw
			defer cw.Close()
		}

		//Подменяем ридер
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		//Основной кот
		h.ServeHTTP(finw, r)
	}

	return http.HandlerFunc(funcCompress)
}
