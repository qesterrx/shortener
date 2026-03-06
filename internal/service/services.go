package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/qesterrx/shortener/internal/model"
)

type URLShortnerStorage struct {
	file    *os.File
	Storage map[string]*model.ShortenURL
}

func GetHash(text string) string {
	hash := sha256.Sum256([]byte(text))
	hex := hex.EncodeToString(hash[:])

	return hex[:8]
}

func NewURLShortnerStorage(filename string) (*URLShortnerStorage, error) {

	storage := make(map[string]*model.ShortenURL)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(file)

	for {
		val := model.ShortenURL{}
		err := dec.Decode(&val)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		storage[val.ShortURL] = &val
	}

	s := URLShortnerStorage{
		file:    file,
		Storage: storage,
	}

	//s.Show()

	return &s, nil
}

func (s *URLShortnerStorage) Get(shortURL string) (string, error) {
	val, exists := s.Storage[shortURL]

	if exists {
		return val.OriginalURL, nil
	}

	return "", fmt.Errorf("object %s not found ", shortURL)

}

func (s *URLShortnerStorage) Set(url string) (string, error) {

	shortURL := GetHash(url)

	if val, exists := s.Storage[shortURL]; exists {
		if val.OriginalURL == url {
			return shortURL, nil
		}
		return "", fmt.Errorf("hash is busy")
	} else {
		newShortenURL := model.ShortenURL{UUID: len(s.Storage) + 1, ShortURL: shortURL, OriginalURL: url}
		json, _ := json.Marshal(&newShortenURL)
		s.file.Write(json)
		s.file.WriteString("\n")

		s.Storage[shortURL] = &newShortenURL
		return shortURL, nil
	}

}

func (s *URLShortnerStorage) Show() {
	fmt.Println(`=========SHOW===========`)
	for k, v := range s.Storage {
		fmt.Printf("%s = %s \n", k, v.OriginalURL)
	}
	fmt.Println(`----------END----------`)
}

func (s *URLShortnerStorage) Close() error {

	return s.file.Close()

}
