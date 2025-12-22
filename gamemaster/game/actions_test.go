package game

import (
	"context"
	"errors"
	"reflect"
	"testing"

	bft "github.com/corstijank/mekstrike/domain/battlefield"
	ut "github.com/corstijank/mekstrike/domain/unit"
	bf "github.com/corstijank/mekstrike/gamemaster/clients/bfclient"
	uc "github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	dclient "github.com/dapr/go-sdk/client"
)

func TestGetAvailableActions_Success(t *testing.T) {
	// Save originals and restore after
	origGetUnit := getUnitClient
	origGetBF := getBattlefieldClient
	t.Cleanup(func() {
		getUnitClient = origGetUnit
		getBattlefieldClient = origGetBF
	})

	// Mock unit client
	getUnitClient = func(_ dclient.Client, id string) *uc.UnitClient {
		return &uc.UnitClient{
			GetData: func(ctx context.Context) (ut.Data, error) {
				return ut.Data{
					Owner:    "Alice",
					Location: &ut.Location{Position: &bft.Coordinates{X: 1, Y: 1}},
					Stats:    &ut.Stats{Movement: "6\""},
				}, nil
			},
		}
	}

	// Mock battlefield client
	coords := []bft.Coordinates{{X: 2, Y: 1}, {X: 2, Y: 2}}
	getBattlefieldClient = func(_ dclient.Client, id string) *bf.BattlefieldClient {
		return &bf.BattlefieldClient{
			GetMovementOptions: func(ctx context.Context, d *ut.Data) ([]bft.Coordinates, error) {
				return coords, nil
			},
		}
	}

	d := Data{
		ID:               "G1",
		CurrentRound:     3,
		CurrentPhase:     Movement,
		CurrentUnitIdx:   0,
		CurrentUnitOrder: []string{"U1"},
		Battlefieldld:    "BF1",
	}

	got, err := d.GetAvailableActions(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.GameID != "G1" || got.CurrentRound != 3 || got.CurrentPhase != Movement || got.CurrentUnitID != "U1" {
		t.Fatalf("options mismatch: %+v", got.Options)
	}
	if got.UnitOwner != "Alice" || got.ActionType != "movement" {
		t.Fatalf("owner/type mismatch: %+v", got)
	}
	if !reflect.DeepEqual(got.AllowedCoordinates, coords) {
		t.Fatalf("coords mismatch: %+v", got.AllowedCoordinates)
	}
}

func TestGetAvailableActions_Errors(t *testing.T) {
	origGetUnit := getUnitClient
	origGetBF := getBattlefieldClient
	t.Cleanup(func() {
		getUnitClient = origGetUnit
		getBattlefieldClient = origGetBF
	})

	// Case 1: unit GetData fails
	getUnitClient = func(_ dclient.Client, id string) *uc.UnitClient {
		return &uc.UnitClient{
			GetData: func(ctx context.Context) (ut.Data, error) { return ut.Data{}, errors.New("boom") },
		}
	}
	d := Data{CurrentUnitOrder: []string{"U1"}}
	if _, err := d.GetAvailableActions(context.Background(), nil); err == nil {
		t.Fatalf("expected error from unit GetData")
	}

	// Case 2: battlefield GetMovementOptions fails
	getUnitClient = func(_ dclient.Client, id string) *uc.UnitClient {
		return &uc.UnitClient{
			GetData: func(ctx context.Context) (ut.Data, error) { return ut.Data{}, nil },
		}
	}
	getBattlefieldClient = func(_ dclient.Client, id string) *bf.BattlefieldClient {
		return &bf.BattlefieldClient{
			GetMovementOptions: func(ctx context.Context, d *ut.Data) ([]bft.Coordinates, error) { return nil, errors.New("nope") },
		}
	}
	if _, err := d.GetAvailableActions(context.Background(), nil); err == nil {
		t.Fatalf("expected error from battlefield GetMovementOptions")
	}
}
