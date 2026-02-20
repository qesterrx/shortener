package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	ServerHost     NetAddress
	ServerRedirect NetAddress
	FileStorage    string
}

type NetAddress struct {
	Host string
	Port int
}

func (t *NetAddress) String() string {
	return t.Host + ":" + strconv.Itoa(t.Port)
}

func (t *NetAddress) Set(val string) error {

	//Тестом на вход подаются почему то http://localhost:8080
	val = strings.ReplaceAll(val, "http://", "")
	val = strings.ReplaceAll(val, "https://", "")

	p := strings.Split(val, ":")
	if len(p) != 2 {
		return errors.New("params must have format host:port")
	}

	port, err := strconv.Atoi(p[1])
	if err != nil {
		return err
	}

	t.Host = p[0]
	t.Port = port
	return nil

}

func ParseParams() *Configuration {
	var cfg Configuration
	cfg.ServerHost = NetAddress{Host: "localhost", Port: 8080}
	flag.Var(&cfg.ServerHost, "a", "Host for request short url")

	cfg.ServerRedirect = NetAddress{Host: "localhost", Port: 8080}
	flag.Var(&cfg.ServerRedirect, "b", "Host for redirect by short url")

	flag.StringVar(&cfg.FileStorage, "f", "TempFileStorage", "Path to file for storage")

	flag.Parse()

	if envServerHost := os.Getenv("SERVER_ADDRESS"); envServerHost != "" {
		cfg.ServerHost.Set(envServerHost)
	}

	if envServerRedirect := os.Getenv("BASE_URL"); envServerRedirect != "" {
		cfg.ServerRedirect.Set(envServerRedirect)
	}

	if envFileStorage := os.Getenv("FILE_STORAGE_PATH"); envFileStorage != "" {
		cfg.FileStorage = envFileStorage
	}

	return &cfg
}
