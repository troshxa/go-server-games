package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*Game, error)
	CreateGame(ctx context.Context, p *Game) error
}

type gameRepo struct {
	db *pgxpool.Pool
}

func NewGameRepo(db *pgxpool.Pool) GameRepo {
	return &gameRepo{db: db}
}

func (r *gameRepo) GetById(ctx context.Context, id uuid.UUID) (*Game, error) {
	query := `SELECT id, title, release_date, price, rating, created_at FROM games WHERE id = $1`
	p := &Game{}
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Title, &p.ReleaseDate, &p.Price, &p.Rating, &p.CreatedAt)
	if err != nil {
		fmt.Print("get game by id: \n", err)

		return nil, err
	}
	return p, nil
}

func (r *gameRepo) CreateGame(ctx context.Context, p *Game) error {
	query := `INSERT INTO games (title, release_date, price, rating) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, p.Title, p.ReleaseDate, p.Price, p.Rating)
	if err != nil {
		fmt.Print("create game: \n", err)
		return err
	}

	return nil
}
