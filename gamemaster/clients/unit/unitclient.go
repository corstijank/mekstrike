package unit

import (
	"context"
	"fmt"

	"github.com/corstijank/mekstrike/src/common/go/unit"
	dapr "github.com/dapr/go-sdk/client"
)

type UnitClient struct {
	daprClient dapr.Client
	id         string
	Deploy     func(context.Context, DeployData) error
}

type DeployData struct {
	Owner          string
	BattlefieldID  string
	Stats          *unit.UnitStats
	DeployLocation string
}

func NewUnit(battlefiedID string, player string, stats *unit.UnitStats) (*UnitClient, error) {
	var err error
	client, err := dapr.NewClient()
	if err != nil {
		return &UnitClient{}, err
	}
	id := fmt.Sprintf("%s-%s-%s", stats.GetModel(), player, battlefiedID)

	result := &UnitClient{
		daprClient: client,
		id:         id,
	}

	client.ImplActorClientStub(result)

	return result, nil
}

func (a *UnitClient) Type() string {
	return "unit"
}

func (a *UnitClient) ID() string {
	return a.id
}
