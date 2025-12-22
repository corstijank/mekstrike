package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	uc "github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	"github.com/corstijank/mekstrike/gamemaster/game"
	"github.com/corstijank/mekstrike/gamemaster/internal/repository"
	"github.com/corstijank/mekstrike/gamemaster/internal/types"
	dclient "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
)

func TestWriteJSONResponse_SetsHeadersAndBody(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := map[string]any{"ok": true}
	writeJSONResponse(payload, rr)

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("CORS header = %q, want *", got)
	}

	var rt map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &rt); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if ok, _ := rt["ok"].(bool); !ok {
		t.Fatalf("response payload mismatch: %v", rt)
	}
}

func TestWriteJSONResponse_ErrorPath(t *testing.T) {
	rr := httptest.NewRecorder()
	// json.Marshal fails on functions -> triggers error path
	payload := map[string]any{"bad": func() {}}
	writeJSONResponse(payload, rr)

	if rr.Code != 500 {
		t.Fatalf("status = %d, want 500", rr.Code)
	}
}

func TestGetGame_Success(t *testing.T) {
	repo := repository.NewMemoryGameRepository()
	g := &game.Data{ID: "G1"}
	_ = repo.Save(nil, g)

	h := NewGameHandlers(nil, repo)

	req := httptest.NewRequest(http.MethodGet, "/games/G1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "G1"})
	rr := httptest.NewRecorder()

	h.GetGame(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var rt game.Data
	if err := json.Unmarshal(rr.Body.Bytes(), &rt); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if rt.ID != "G1" {
		t.Fatalf("ID = %s, want G1", rt.ID)
	}
}

func TestGetGame_NotFound(t *testing.T) {
	repo := repository.NewMemoryGameRepository()
	h := NewGameHandlers(nil, repo)

	req := httptest.NewRequest(http.MethodGet, "/games/NOPE", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "NOPE"})
	rr := httptest.NewRecorder()

	h.GetGame(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}

func TestGetGames_ReturnsSortedList(t *testing.T) {
	repo := repository.NewMemoryGameRepository()
	_ = repo.Save(nil, &game.Data{ID: "G2"})
	_ = repo.Save(nil, &game.Data{ID: "G1"})
	h := NewGameHandlers(nil, repo)

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	rr := httptest.NewRecorder()
	h.GetGames(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var rt []game.Data
	if err := json.Unmarshal(rr.Body.Bytes(), &rt); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(rt) != 2 {
		t.Fatalf("len = %d, want 2", len(rt))
	}
}

func TestMoveUnit_ValidatesTurnOwnership(t *testing.T) {
	repo := repository.NewMemoryGameRepository()
	// Current unit is U1, request will use U2 -> should 400
	_ = repo.Save(nil, &game.Data{ID: "G1", CurrentUnitOrder: []string{"U1"}, CurrentUnitIdx: 0})
	h := NewGameHandlers(nil, repo)

	body, _ := json.Marshal(types.MoveRequest{UnitId: "U2", X: 0, Y: 0, Heading: 0})
	req := httptest.NewRequest(http.MethodPost, "/games/G1/move", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "G1"})
	rr := httptest.NewRecorder()
	h.MoveUnit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}

func TestMoveUnit_HappyPath(t *testing.T) {
	// Mock the unit client used by the handler
	orig := getUnitClientInHandler
	t.Cleanup(func() { getUnitClientInHandler = orig })

	var captured map[string]interface{}
	getUnitClientInHandler = func(_ dclient.Client, id string) *uc.UnitClient {
		return &uc.UnitClient{
			Move: func(_ context.Context, req interface{}) error {
				// Handler builds a map[string]interface{}
				if m, ok := req.(map[string]interface{}); ok {
					captured = m
				}
				return nil
			},
		}
	}

	repo := repository.NewMemoryGameRepository()
	_ = repo.Save(nil, &game.Data{ID: "G1", CurrentUnitOrder: []string{"U1"}, CurrentUnitIdx: 0})
	h := NewGameHandlers(nil, repo)

	body, _ := json.Marshal(types.MoveRequest{UnitId: "U1", X: 5, Y: 6, Heading: 3})
	req := httptest.NewRequest(http.MethodPost, "/games/G1/move", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "G1"})
	rr := httptest.NewRecorder()

	h.MoveUnit(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["success"] != true || resp["unitId"] != "U1" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	// Verify captured move request shape (numbers may be int32)
	toInt := func(v interface{}) int {
		switch vv := v.(type) {
		case int:
			return vv
		case int32:
			return int(vv)
		case int64:
			return int(vv)
		case float64:
			return int(vv)
		default:
			return -1
		}
	}
	if captured == nil || toInt(captured["x"]) != 5 || toInt(captured["y"]) != 6 || toInt(captured["heading"]) != 3 {
		t.Fatalf("captured move mismatch: %+v", captured)
	}
}

func TestGetGameLogs_TransformsCloudEvents(t *testing.T) {
	repo := repository.NewMemoryGameRepository()
	_ = repo.Save(nil, &game.Data{ID: "G1", GameLogs: []types.CloudEvent{{Type: "t1", Data: "m1", Time: "2024-01-01T00:00:00Z"}}})
	h := NewGameHandlers(nil, repo)

	req := httptest.NewRequest(http.MethodGet, "/games/G1/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "G1"})
	rr := httptest.NewRecorder()
	h.GetGameLogs(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var rt []map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &rt); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(rt) != 1 || rt[0]["type"] != "t1" || rt[0]["message"] != "m1" {
		t.Fatalf("unexpected payload: %+v", rt)
	}
}
