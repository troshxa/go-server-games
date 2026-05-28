package main

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	Price       float64   `json:"price"`
	Rating      int       `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

type GamePatch struct {
	Title       *string    `json:"title"`
	ReleaseDate *time.Time `json:"release_date"`
	Price       *float64   `json:"price"`
	Rating      *int       `json:"rating"`
}
