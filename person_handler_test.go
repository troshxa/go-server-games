package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
