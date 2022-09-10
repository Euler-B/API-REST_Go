package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Euler-B/API-REST_Go/models"
	"github.com/Euler-B/API-REST_Go/repository"
	"github.com/Euler-B/API-REST_Go/server"
)

const(
	HASH_COST = 8
)
type SignUpLoginRequest struct {
	Email      string `json:"Email"`
	Password   string `json:"password"`
}

type SignUpResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
}

type LoginResponse struct {
	Token    string `json:"string"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 
		HASH_COST)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var user = models.User{
			Email:     request.Email,
			Password:  string(hashedPassword),
			Id :       id.String(),
		}
		err = repository.InsertUser(r.Context(), &user)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:       user.Id,
			Email:    user.Email,
		})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request)  {
		var request = SignUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "INVALID CREDENTIALS", http.StatusUnauthorized)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), 
		[]byte(request.Password)); err != nil {
			http.Error(w, "INVALID CREDENTIALS", http.StatusUnauthorized)
			return
		}
		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				 ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})
	}
}

