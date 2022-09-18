package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"

	"github.com/Euler-B/API-REST_Go/models"
	"github.com/Euler-B/API-REST_Go/repository"
	"github.com/Euler-B/API-REST_Go/server"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

type PostUpdateResponse struct {
	Message     string  `json:"message"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, 
			func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
				var postRequest = UpsertPostRequest{}
				err := json.NewDecoder(r.Body).Decode(&postRequest); 
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				id, err := ksuid.NewRandom()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				post := models.Post{
					Id: id.String(),
					PostContent: postRequest.PostContent,
					UserId: claims.UserId,
				}
				err = repository.InsertPost(r.Context(), &post)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(PostResponse{
					Id: post.Id,
					PostContent: post.PostContent,
				})
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	}
}

func GetPostByIdHandler(s server.Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		params    := mux.Vars(r) 
		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params    := mux.Vars(r)
		tokenString := strings.TrimSpace(r.Header.Get("Authorization")) // reto propuesto:
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, // buscar una forma de no repetir tanto el mismo codigo para validar tokens
			func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
				var postRequest = UpsertPostRequest{}
				err := json.NewDecoder(r.Body).Decode(&postRequest); 
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				post := models.Post{
					Id: params["id"],
					PostContent: postRequest.PostContent,
					UserId: claims.UserId,
				}
				err = repository.UpdatePost(r.Context(), &post)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(PostUpdateResponse{
					Message: "Post Updated",
				})
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	}
}

func DeletePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params    := mux.Vars(r)
		tokenString := strings.TrimSpace(r.Header.Get("Authorization")) 
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
				err = repository.DeletePost(r.Context(), params["id"], claims.UserId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(PostUpdateResponse{
					Message: "Post Deleted",
				})
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	}
}