package game

import (
	"context"
	"encoding/json"
	"log"

	dapr "github.com/dapr/go-sdk/client"
)

func (g *Data) isAIUnit(unitID string) bool {
	// Check if unit belongs to CPU player
	for _, cpuUnit := range g.PlayerBUnits {
		if cpuUnit == unitID {
			return true
		}
	}
	return false
}

func (g *Data) publishAITurnEvent(ctx context.Context, client dapr.Client, unitID string) error {
	eventData := map[string]interface{}{
		"gameId": g.ID,
		"unitId": unitID,
		"phase":  g.phaseToString(),
		"round":  g.CurrentRound,
	}

	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	err = client.PublishEvent(ctx, "redis-pubsub", "ai-turn-started", eventJSON)
	if err != nil {
		return err
	}

	log.Printf("Published AI turn event: %s", string(eventJSON))
	return nil
}

func (g *Data) phaseToString() string {
	switch g.CurrentPhase {
	case Movement:
		return "Movement"
	case Combat:
		return "Combat"
	case End:
		return "End"
	default:
		return "Unknown"
	}
}