package bfclient

import (
	"context"

	"github.com/corstijank/mekstrike/domain/battlefield"
	unitType "github.com/corstijank/mekstrike/domain/unit"
	dapr "github.com/dapr/go-sdk/client"
)

type BattlefieldClient struct {
	id                 string
	GetNumberOfCols    func(context.Context) (int, error)
	GetNumberOfRows    func(context.Context) (int, error)
	GetBoardCells      func(context.Context) ([]battlefield.Cell, error)
	GetMovementOptions func(context.Context, *unitType.Data) ([]battlefield.Coordinates, error)
}

func GetBattlefieldClient(client dapr.Client, id string) *BattlefieldClient {
	result := &BattlefieldClient{
		id: id,
	}

	client.ImplActorClientStub(result)
	return result
}

func (a *BattlefieldClient) Type() string {
	return "battlefield"
}

func (a *BattlefieldClient) ID() string {
	return a.id
}
