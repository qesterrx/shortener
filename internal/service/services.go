package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type URLShortnerStorage struct {
	Storage map[string]string
}

func (s *URLShortnerStorage) Get(shortURL string) (string, error) {
	val, exists := s.Storage[shortURL]

	if exists {
		return val, nil
	}

	return "", fmt.Errorf("object %s not found ", shortURL)

}

func (s *URLShortnerStorage) Set(url string) (string, error) {

	hash := sha256.Sum256([]byte(url))
	hex := hex.EncodeToString(hash[:])
	shortURL := hex[:8]

	if val, exists := s.Storage[shortURL]; exists {
		if val == url {
			return shortURL, nil
		}
		return "", fmt.Errorf("hash is busy")
	} else {
		s.Storage[shortURL] = url
		return shortURL, nil
	}

}

func (s *URLShortnerStorage) Show() {
	fmt.Println(`=========SHOW===========`)
	for k, v := range s.Storage {
		fmt.Printf("%s = %s \n", k, v)
	}
	fmt.Println(`----------END----------`)
}
