package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/segmentio/ksuid"

	"github.com/Euler-B/API-REST_Go/models"
	"github.com/Euler-B/API-REST_Go/repository"
	"github.com/Euler-B/API-REST_Go/server"
)

type SignUpRequest struct {
	Email    string `json:"Email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var user = models.User{
			Email:     request.Email,
			Password:  request.Password,
			Id :       id.String(),
		}
		err = repository.InsertUser(r.Context(), &user)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:     user.Id,
			Email:  user.Email,
		})
	}
}

