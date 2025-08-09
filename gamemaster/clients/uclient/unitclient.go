package uclient

import (
	"context"
	"fmt"

	"github.com/corstijank/mekstrike/domain/unit"
	dapr "github.com/dapr/go-sdk/client"
)

type UnitClient struct {
	id        string
	Deploy    func(context.Context, unit.DeployRequest) error
	SetActive func(context.Context, bool) error
	GetData   func(context.Context) (unit.Data, error)
	Move      func(context.Context, interface{}) error
}

func NewUnit(client dapr.Client, battlefiedID string, player string, stats *unit.Stats) *UnitClient {
	result := &UnitClient{
		id: fmt.Sprintf("%s-%s-%s", stats.GetModel(), player, battlefiedID),
	}
	client.ImplActorClientStub(result)

	return result
}

func GetUnitClient(client dapr.Client, id string) *UnitClient {
	result := &UnitClient{
		id: id,
	}
	client.ImplActorClientStub(result)

	return result
}

func (a *UnitClient) Type() string {
	return "unit"
}

func (a *UnitClient) ID() string {
	return a.id
}
