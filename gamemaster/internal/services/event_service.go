package services

import (
	"context"
	"log"

	"github.com/corstijank/mekstrike/gamemaster/game"
	"github.com/corstijank/mekstrike/gamemaster/internal/repository"
	"github.com/corstijank/mekstrike/gamemaster/internal/types"
	dapr "github.com/dapr/go-sdk/client"
)

type EventService struct {
	client     dapr.Client
	repository repository.GameRepository
}

func NewEventService(client dapr.Client, repo repository.GameRepository) *EventService {
	return &EventService{
		client:     client,
		repository: repo,
	}
}

func (s *EventService) ProcessPhaseCompletion(ctx context.Context, event *types.ActionCompletedEvent, expectedPhase string) error {
	log.Printf("Processing %s completion event: %+v", expectedPhase, *event)
	
	gameData, err := s.repository.Get(ctx, event.GameId)
	if err != nil {
		log.Printf("No game found with ID %s for %s completion", event.GameId, expectedPhase)
		return nil // Return nil to avoid retries for missing games
	}
	
	// Idempotency check: verify this event matches current game state
	if !s.isValidPhaseEvent(gameData, event.UnitId, event.Phase, expectedPhase) {
		log.Printf("Ignoring duplicate or out-of-order %s event for game %s, unit %s", expectedPhase, event.GameId, event.UnitId)
		return nil // Return nil to avoid retries for duplicate events
	}
	
	gameData.AdvanceTurn(ctx, s.client)
	
	err = s.repository.Save(ctx, gameData)
	if err != nil {
		log.Printf("Error saving game after %s completion: %v", expectedPhase, err)
		return err
	}
	
	log.Printf("Advanced turn for game %s after %s completion", event.GameId, expectedPhase)
	return nil
}

func (s *EventService) isValidPhaseEvent(gameData *game.Data, unitId, eventPhase, expectedPhase string) bool {
	// Check if we're in the expected phase
	currentPhaseName := s.getPhaseString(gameData.CurrentPhase)
	if currentPhaseName != expectedPhase {
		return false
	}
	
	// Check if it's the right unit's turn
	if gameData.CurrentUnitIdx >= len(gameData.CurrentUnitOrder) {
		return false
	}
	
	currentUnitId := gameData.CurrentUnitOrder[gameData.CurrentUnitIdx]
	if currentUnitId != unitId {
		return false
	}
	
	// Verify the phase in the event matches expected phase
	return eventPhase == expectedPhase
}

func (s *EventService) getPhaseString(phase game.Phase) string {
	switch phase {
	case game.Movement:
		return "Movement"
	case game.Combat:
		return "Combat"
	case game.End:
		return "End"
	default:
		return "Unknown"
	}
}