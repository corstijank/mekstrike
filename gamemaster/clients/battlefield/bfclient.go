package battlefield

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"
)

type BattlefieldClient struct {
	daprClient    dapr.Client
	id            string
	GetBoardCells func(context.Context) ([]Cell, error)
}

type Cell struct {
	Col           int
	Row           int
	TerrainTypeID int
}

func NewBattlefield(id string) (*BattlefieldClient, error) {
	var err error
	client, err := dapr.NewClient()
	if err != nil {
		return &BattlefieldClient{}, err
	}

	result := &BattlefieldClient{
		daprClient: client,
		id:         id,
	}

	client.ImplActorClientStub(result)
	return result, nil
}

func (a *BattlefieldClient) Type() string {
	return "battlefield"
}

func (a *BattlefieldClient) ID() string {
	return a.id
}
