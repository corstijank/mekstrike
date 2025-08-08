package game

import (
	"context"

	"github.com/corstijank/mekstrike/gamemaster/clients/bfclient"
	"github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	dapr "github.com/dapr/go-sdk/client"
)

func (d Data) GetAvailableActions(ctx context.Context, client dapr.Client) (AvailableActions, error) {
	currentUnitID := d.CurrentUnitOrder[d.CurrentUnitIdx]
	currentUnit := uclient.GetUnitClient(client, currentUnitID)
	unitData, err := currentUnit.GetData(ctx)
	if err != nil {
		return AvailableActions{}, err
	}

	bf := bfclient.GetBattlefieldClient(client, d.Battlefieldld)
	options, err := bf.GetMovementOptions(ctx, &unitData)
	if err != nil {
		return AvailableActions{}, err
	}

	return AvailableActions{
		Options: Options{
			GameID:        d.ID,
			CurrentRound:  d.CurrentRound,
			CurrentPhase:  d.CurrentPhase,
			CurrentUnitID: currentUnitID,
		},
		UnitOwner:          unitData.Owner,
		ActionType:         "movement",
		AllowedCoordinates: options}, nil
}