package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type GameHandler struct {
	repo GameRepo
}

func newGameHandler(repo GameRepo) *GameHandler {
	return &GameHandler{repo: repo}
}

func (h *GameHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	p, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(p)
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var p Game
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(p.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if p.Price < 0 {
		http.Error(w, "Price must be a positive number", http.StatusBadRequest)
		return
	}

	err = h.repo.CreateGame(r.Context(), &p)
	if err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(p)
}
