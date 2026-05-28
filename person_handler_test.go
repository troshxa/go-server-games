package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type mockGameRepo struct {
	games map[uuid.UUID]*Game
}

func newMockGameRepo() *mockGameRepo {
	return &mockGameRepo{games: make(map[uuid.UUID]*Game)}
}

func (m *mockGameRepo) CreateGame(ctx context.Context, p *Game) error {
	m.games[p.ID] = p
	return nil
}

func (m *mockGameRepo) GetById(ctx context.Context, id uuid.UUID) (*Game, error) {
	p, ok := m.games[id]
	if !ok {
		return nil, fmt.Errorf("game not found")
	}

	return p, nil
}

func TestCreateGameHandler_InvalidJson(t *testing.T) {
	repo := newMockGameRepo()
	handler := newGameHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/games", strings.NewReader("{dasdasd}"))
	rec := httptest.NewRecorder()

	handler.CreateGame(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateGameHandler_EmptyTitle(t *testing.T) {
	repo := newMockGameRepo()
	handler := newGameHandler(repo)

	body := `{"title":"","release_date":"2020-01-01T00:00:00Z","price":29.99,"rating":8}`
	req := httptest.NewRequest(http.MethodPost, "/games", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.CreateGame(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateGameHandler_NegativePrice(t *testing.T) {
	repo := newMockGameRepo()
	handler := newGameHandler(repo)

	body := `{"title":"Witcher","release_date":"2020-01-01T00:00:00Z","price":-5,"rating":8}`
	req := httptest.NewRequest(http.MethodPost, "/games", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.CreateGame(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateGameHandler_InvalidRating(t *testing.T) {
	repo := newMockGameRepo()
	handler := newGameHandler(repo)

	body := `{"title":"Witcher","release_date":"2020-01-01T00:00:00Z","price":29.99,"rating":15}`
	req := httptest.NewRequest(http.MethodPost, "/games", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.CreateGame(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateGameHandler_Success(t *testing.T) {
	repo := newMockGameRepo()
	handler := newGameHandler(repo)

	body := `{"title":"Witcher","release_date":"2020-01-01T00:00:00Z","price":29.99,"rating":9}`
	req := httptest.NewRequest(http.MethodPost, "/games", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.CreateGame(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var result Game
	err := json.NewDecoder(rec.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Title != "Witcher" {
		t.Errorf("expected title %q, got %q", "Witcher", result.Title)
	}
}

func chiCtxWithID(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/games/"+id, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestGetByIdHandler_InvalidUUID(t *testing.T) {
	repo := newMockGameRepo()
	handler := newGameHandler(repo)

	req := chiCtxWithID("not-a-uuid")
	rec := httptest.NewRecorder()

	handler.GetById(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetByIdHandler_NotFound(t *testing.T) {
	repo := newMockGameRepo()
	handler := newGameHandler(repo)

	req := chiCtxWithID(uuid.New().String())
	rec := httptest.NewRecorder()

	handler.GetById(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetByIdHandler_Success(t *testing.T) {
	repo := newMockGameRepo()
	handler := newGameHandler(repo)

	id := uuid.New()
	repo.games[id] = &Game{
		ID:          id,
		Title:       "Celeste",
		ReleaseDate: time.Date(2018, 1, 25, 0, 0, 0, 0, time.UTC),
		Price:       19.99,
		Rating:      9,
	}

	req := chiCtxWithID(id.String())
	rec := httptest.NewRecorder()

	handler.GetById(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var result Game
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.ID != id {
		t.Errorf("expected id %v, got %v", id, result.ID)
	}
	if result.Title != "Celeste" {
		t.Errorf("expected title %q, got %q", "Celeste", result.Title)
	}
}
