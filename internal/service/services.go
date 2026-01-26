package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type URLShortnerStorage struct {
	Storage map[string]string
}

func (s *URLShortnerStorage) Get(shortUrl string) string {
	val, _ := s.Storage[shortUrl]
	return val
}

func (s *URLShortnerStorage) Set(url string) string {

	hash := sha256.Sum256([]byte(url))
	hex := hex.EncodeToString(hash[:])
	shortUrl := hex[:8]

	if val, exists := s.Storage[shortUrl]; exists {
		if val == url {
			return shortUrl
		}
		return s.Set(url + hex)
	} else {
		s.Storage[shortUrl] = url
		return shortUrl
	}

}

func (s *URLShortnerStorage) Show() {
	fmt.Println(`=========SHOW===========`)
	for k, v := range s.Storage {
		fmt.Printf("%s = %s \n", k, v)
	}
	fmt.Println(`----------END----------`)
}
