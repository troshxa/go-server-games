package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("dotenv")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	dsn := viper.GetString("DATABASE_URL")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is not set")
		os.Exit(1)
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	if err := RunMigrations(dsn); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run migrations: %v\n", err)
		os.Exit(1)
	}
	log.Println("Migrations applied")

	dbpool, err := pgxpool.New(context.Background(), dsn)
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

	log.Fatal(http.ListenAndServe(":"+port, r))
}
