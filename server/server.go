package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Euler-B/API-REST_Go/repository"
	"github.com/Euler-B/API-REST_Go/database"
	"github.com/Euler-B/API-REST_Go/websocket"

)

type Config struct {
	Port string
	JWTSecret string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub()    *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub

}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("el puerto es requerido")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("secret es requerido")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("la base de datos es requerida")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}

	return broker, nil
}

func (b *Broker) Start (binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	
	repo, err := database.NewPostgreRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)

	log.Println("Inicializando servidor en el Puerto", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}