package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Euler-B/API-REST_Go/handlers"
	"github.com/Euler-B/API-REST_Go/middleware"
	"github.com/Euler-B/API-REST_Go/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	PORT           := os.Getenv("PORT")
	JWT_SECRET     := os.Getenv("JWT_SECRET")
	DATABASE_URL   := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config {
		JWTSecret:    JWT_SECRET,
		Port:         PORT,
		DatabaseUrl:  DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {
	r.Use(middleware.CheckAuthMiddleware(s))

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	
	r.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
}