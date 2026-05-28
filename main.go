package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/some-db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	gameRepo := NewGameRepo(dbpool)
	gameHandler := newGameHandler(gameRepo)

	r := chi.NewRouter()

	r.Get("/games", gameHandler.GetAll)
	r.Get("/games/{id}", gameHandler.GetById)
	r.Post("/games", gameHandler.CreateGame)
	r.Put("/games/{id}", gameHandler.UpdateGame)
	r.Patch("/games/{id}", gameHandler.PatchGame)
	r.Delete("/games/{id}", gameHandler.DeleteGame)

	log.Fatal(http.ListenAndServe(":8080", r))
}
