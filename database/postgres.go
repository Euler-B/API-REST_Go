package database

import (
	"context"
	"database/sql"

	"github.com/Euler-B/API-REST_Go/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgreREpository(ulr string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", ulr)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
} 

func (repo *PostgresRepository) InsertUser (ctx context.Context, user *models.User) error {
	repo.db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2, $3)") 
	
}