package game

import (
	"context"
	"log"

	"github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	dapr "github.com/dapr/go-sdk/client"
)

func (g *Data) NewRound(ctx context.Context, client dapr.Client) {
	g.CurrentRound++
	g.CurrentPhase = Movement
	g.CurrentUnitIdx = 0
	log.Printf("New Round: %d", g.CurrentRound)
	g.activateCurrentUnit(ctx, client)
}

func (g *Data) AdvanceTurn(ctx context.Context, client dapr.Client) {
	// Deactivate current unit
	if g.CurrentUnitIdx < len(g.CurrentUnitOrder) {
		currentUnitID := g.CurrentUnitOrder[g.CurrentUnitIdx]
		u := uclient.GetUnitClient(client, currentUnitID)
		err := u.SetActive(ctx, false)
		if err != nil {
			log.Printf("Error deactivating unit %s: %v", currentUnitID, err)
		}
	}

	// Advance to next unit
	g.CurrentUnitIdx++
	
	// Check if all units have moved in current phase
	if g.CurrentUnitIdx >= len(g.CurrentUnitOrder) {
		g.advancePhase(ctx, client)
	} else {
		g.activateCurrentUnit(ctx, client)
	}
}

func (g *Data) advancePhase(ctx context.Context, client dapr.Client) {
	switch g.CurrentPhase {
	case Movement:
		log.Printf("Advancing to Combat phase")
		g.CurrentPhase = Combat
		g.CurrentUnitIdx = 0
		g.activateCurrentUnit(ctx, client)
	case Combat:
		log.Printf("Advancing to End phase")
		g.CurrentPhase = End
		g.CurrentUnitIdx = 0
		g.activateCurrentUnit(ctx, client)
	case End:
		log.Printf("Round completed, starting new round")
		g.NewRound(ctx, client)
	}
}

func (g *Data) activateCurrentUnit(ctx context.Context, client dapr.Client) {
	if g.CurrentUnitIdx >= len(g.CurrentUnitOrder) {
		log.Printf("No more units in order, advancing phase or round")
		return
	}

	currentUnitID := g.CurrentUnitOrder[g.CurrentUnitIdx]
	log.Printf("Activating unit %s", currentUnitID)
	
	u := uclient.GetUnitClient(client, currentUnitID)
	err := u.SetActive(ctx, true)
	if err != nil {
		log.Printf("Error setting active\n%v", err.Error())
		return
	}

	// Check if current unit is AI-owned and publish event
	if g.isAIUnit(currentUnitID) {
		err := g.publishAITurnEvent(ctx, client, currentUnitID)
		if err != nil {
			log.Printf("Error publishing AI turn event: %v", err)
		}
	}
}