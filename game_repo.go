package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameRepo interface {
	GetAll(ctx context.Context) ([]*Game, error)
	GetById(ctx context.Context, id uuid.UUID) (*Game, error)
	CreateGame(ctx context.Context, p *Game) error
	UpdateGame(ctx context.Context, id uuid.UUID, p *Game) (*Game, error)
	PatchGame(ctx context.Context, id uuid.UUID, p *GamePatch) (*Game, error)
	DeleteGame(ctx context.Context, id uuid.UUID) (*Game, error)
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

func (r *gameRepo) GetAll(ctx context.Context) ([]*Game, error) {
	query := `SELECT id, title, release_date, price, rating, created_at FROM games`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		fmt.Print("get all games: \n", err)
		return nil, err
	}
	defer rows.Close()

	var games []*Game
	for rows.Next() {
		p := &Game{}
		if err := rows.Scan(&p.ID, &p.Title, &p.ReleaseDate, &p.Price, &p.Rating, &p.CreatedAt); err != nil {
			return nil, err
		}
		games = append(games, p)
	}
	return games, nil
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

func (r *gameRepo) UpdateGame(ctx context.Context, id uuid.UUID, p *Game) (*Game, error) {
	query := `UPDATE games SET title=$1, release_date=$2, price=$3, rating=$4
		WHERE id=$5 RETURNING id, title, release_date, price, rating, created_at`
	result := &Game{}
	err := r.db.QueryRow(ctx, query, p.Title, p.ReleaseDate, p.Price, p.Rating, id).
		Scan(&result.ID, &result.Title, &result.ReleaseDate, &result.Price, &result.Rating, &result.CreatedAt)
	if err != nil {
		fmt.Print("update game: \n", err)
		return nil, err
	}
	return result, nil
}

func (r *gameRepo) PatchGame(ctx context.Context, id uuid.UUID, p *GamePatch) (*Game, error) {
	setClauses := []string{}
	args := []any{}
	i := 1

	if p.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title=$%d", i))
		args = append(args, *p.Title)
		i++
	}
	if p.ReleaseDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("release_date=$%d", i))
		args = append(args, *p.ReleaseDate)
		i++
	}
	if p.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price=$%d", i))
		args = append(args, *p.Price)
		i++
	}
	if p.Rating != nil {
		setClauses = append(setClauses, fmt.Sprintf("rating=$%d", i))
		args = append(args, *p.Rating)
		i++
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE games SET %s WHERE id=$%d
		RETURNING id, title, release_date, price, rating, created_at`,
		strings.Join(setClauses, ", "), i)

	result := &Game{}
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&result.ID, &result.Title, &result.ReleaseDate, &result.Price, &result.Rating, &result.CreatedAt)
	if err != nil {
		fmt.Print("patch game: \n", err)
		return nil, err
	}
	return result, nil
}

func (r *gameRepo) DeleteGame(ctx context.Context, id uuid.UUID) (*Game, error) {
	query := `DELETE FROM games WHERE id=$1
		RETURNING id, title, release_date, price, rating, created_at`
	result := &Game{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&result.ID, &result.Title, &result.ReleaseDate, &result.Price, &result.Rating, &result.CreatedAt)
	if err != nil {
		fmt.Print("delete game: \n", err)
		return nil, err
	}
	return result, nil
}
