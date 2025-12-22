package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	bft "github.com/corstijank/mekstrike/domain/battlefield"
	ut "github.com/corstijank/mekstrike/domain/unit"
	bfclient "github.com/corstijank/mekstrike/gamemaster/clients/bfclient"
	uc "github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	"github.com/corstijank/mekstrike/gamemaster/game"
	"github.com/corstijank/mekstrike/gamemaster/internal/repository"
	dclient "github.com/dapr/go-sdk/client"
)

func TestNewGame_HappyPath(t *testing.T) {
	// Armybuilder returns two units each call
	restoreAB := game.UseArmyBuilderStubForTest([]*ut.Stats{{Model: "FAKE-1"}, {Model: "FAKE-2"}})
	defer restoreAB()

	// Battlefield client used in game.New: only GetBoardCells is called for init; return empty slice
	restoreBF := game.SetBattlefieldClientFactoryForTest(func(_ dclient.Client, id string) *bfclient.BattlefieldClient {
		return &bfclient.BattlefieldClient{
			GetBoardCells: func(ctx context.Context) ([]bft.Cell, error) { return []bft.Cell{}, nil },
		}
	})
	defer restoreBF()

	// Unit client factory: return a client with predictable ID and no-op Deploy
	var deployed []string
	restoreUnit := game.SetUnitClientFactoryForTest(func(_ dclient.Client, bfID, owner string, stats *ut.Stats) *uc.UnitClient {
		id := stats.GetModel() + "-" + owner + "-" + bfID
		return &uc.UnitClient{
			Deploy: func(ctx context.Context, req ut.DeployRequest) error {
				deployed = append(deployed, id)
				return nil
			},
			SetActive: func(ctx context.Context, b bool) error { return nil },
			GetData:   func(ctx context.Context) (ut.Data, error) { return ut.Data{}, nil },
			Move:      func(ctx context.Context, r interface{}) error { return nil },
		}
	})
	defer restoreUnit()

	repo := repository.NewMemoryGameRepository()
	h := NewGameHandlers(nil, repo)

	body, _ := json.Marshal(map[string]string{"playerName": "Alice"})
	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.NewGame(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}

	var resp game.Data
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.PlayerA != "Alice" || resp.Battlefieldld == "" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(resp.PlayerAUnits) == 0 || len(resp.PlayerBUnits) == 0 {
		t.Fatalf("units not deployed as expected")
	}
	if len(deployed) == 0 {
		t.Fatalf("expected deployments to be recorded")
	}
}
