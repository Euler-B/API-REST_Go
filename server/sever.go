package server

import (
	"context"
	"errors"

	"github.com/gorilla/mux"
)

type Config struct {
	Port string
	JWTSecret string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router

}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("El Puerto es Requerido")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("El Secret es Requerido")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("La base de datos es Requerida")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	return broker, nil
}