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

func (h *GameHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	games, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
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
	if p.Rating < 0 || p.Rating > 10 {
		http.Error(w, "Rating must be between 0 and 10", http.StatusBadRequest)
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

func (h *GameHandler) UpdateGame(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p Game
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
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
	if p.Rating < 0 || p.Rating > 10 {
		http.Error(w, "Rating must be between 0 and 10", http.StatusBadRequest)
		return
	}

	result, err := h.repo.UpdateGame(r.Context(), id, &p)
	if err != nil {
		http.Error(w, "Failed to update game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *GameHandler) PatchGame(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p GamePatch
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if p.Title != nil && strings.TrimSpace(*p.Title) == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	if p.Price != nil && *p.Price < 0 {
		http.Error(w, "Price must be a positive number", http.StatusBadRequest)
		return
	}
	if p.Rating != nil && (*p.Rating < 0 || *p.Rating > 10) {
		http.Error(w, "Rating must be between 0 and 10", http.StatusBadRequest)
		return
	}

	result, err := h.repo.PatchGame(r.Context(), id, &p)
	if err != nil {
		http.Error(w, "Failed to patch game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *GameHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := h.repo.DeleteGame(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
