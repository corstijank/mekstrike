package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/corstijank/mekstrike/gamemaster/game"
	"github.com/corstijank/mekstrike/gamemaster/internal/config"
	"github.com/corstijank/mekstrike/gamemaster/internal/repository"
	"github.com/corstijank/mekstrike/gamemaster/internal/types"
	dapr "github.com/dapr/go-sdk/client"
)

type EventHandlers struct {
	client     dapr.Client
	repository repository.GameRepository
	config     *config.Config
}

func NewEventHandlers(client dapr.Client, repo repository.GameRepository, cfg *config.Config) *EventHandlers {
	return &EventHandlers{
		client:     client,
		repository: repo,
		config:     cfg,
	}
}

func (h *EventHandlers) GetDaprSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions := []types.DaprSubscription{
		{
			PubsubName: h.config.PubsubName,
			Topic:      "unit-movement-completed",
			Route:      "unit-movement-completed",
		},
		{
			PubsubName: h.config.PubsubName,
			Topic:      "unit-attack-completed",
			Route:      "unit-attack-completed",
		},
		{
			PubsubName: h.config.PubsubName,
			Topic:      "unit-end-phase-completed",
			Route:      "unit-end-phase-completed",
		},
	}
	writeJSONResponse(subscriptions, w)
}

func (h *EventHandlers) HandleMovementCompleted(w http.ResponseWriter, r *http.Request) {
	h.handlePhaseCompleted(w, r, "Movement")
}

func (h *EventHandlers) HandleAttackCompleted(w http.ResponseWriter, r *http.Request) {
	h.handlePhaseCompleted(w, r, "Combat")
}

func (h *EventHandlers) HandleEndPhaseCompleted(w http.ResponseWriter, r *http.Request) {
	h.handlePhaseCompleted(w, r, "End")
}

func (h *EventHandlers) handlePhaseCompleted(w http.ResponseWriter, r *http.Request, expectedPhase string) {
	event, cloudEvent, err := parseActionCompletedEvent(r)
	if err != nil {
		log.Printf("Error parsing %s completed event: %v", expectedPhase, err)
		http.Error(w, "Invalid event data", 400)
		return
	}
	
	log.Printf("Received %s completed event: %+v", expectedPhase, *event)
	
	gameData, err := h.repository.Get(r.Context(), event.GameId)
	if err != nil {
		log.Printf("No game found with ID %s for %s completion", event.GameId, expectedPhase)
		w.WriteHeader(200) // Still return success to avoid retries
		return
	}
	
	// Idempotency check: verify this event matches current game state
	if !isValidPhaseEvent(gameData, event.UnitId, event.Phase, expectedPhase) {
		log.Printf("Ignoring duplicate or out-of-order %s event for game %s, unit %s", expectedPhase, event.GameId, event.UnitId)
		w.WriteHeader(200) // Still return success to avoid retries
		return
	}
	
	// Store the CloudEvent in game logs if it was a proper CloudEvent
	if cloudEvent != nil {
		gameData.GameLogs = append(gameData.GameLogs, *cloudEvent)
	}
	
	gameData.AdvanceTurn(r.Context(), h.client)
	
	err = h.repository.Save(r.Context(), gameData)
	if err != nil {
		log.Printf("Error saving game after %s completion: %v", expectedPhase, err)
	}
	
	log.Printf("Advanced turn for game %s after %s completion", event.GameId, expectedPhase)
	w.WriteHeader(200)
}

func parseActionCompletedEvent(r *http.Request) (*types.ActionCompletedEvent, *types.CloudEvent, error) {
	// Try to parse as CloudEvent first
	var cloudEvent types.CloudEvent
	body := make([]byte, 0)
	if data, err := io.ReadAll(r.Body); err == nil {
		body = data
		r.Body = io.NopCloser(bytes.NewReader(body)) // Reset body for potential re-read
	}

	if err := json.Unmarshal(body, &cloudEvent); err == nil && cloudEvent.Data != "" {
		// Successfully parsed as CloudEvent, now parse the nested data
		log.Printf("Parsing CloudEvent: data=%s, topic=%s, type=%s", cloudEvent.Data, cloudEvent.Topic, cloudEvent.Type)
		
		var event types.ActionCompletedEvent
		if err := json.Unmarshal([]byte(cloudEvent.Data), &event); err != nil {
			return nil, nil, fmt.Errorf("failed to parse CloudEvent data field: %v", err)
		}
		
		log.Printf("Parsed CloudEvent data: %+v", event)
		return &event, &cloudEvent, nil
	}

	// Fall back to direct parsing (backward compatibility)
	log.Printf("Trying direct parsing of body: %s", string(body))
	var event types.ActionCompletedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, nil, fmt.Errorf("failed to parse as both CloudEvent and direct ActionCompletedEvent: %v", err)
	}
	
	log.Printf("Parsed direct event: %+v", event)
	return &event, nil, nil
}

// isValidPhaseEvent checks if the received event matches the current game state
// This prevents duplicate processing of the same phase completion
func isValidPhaseEvent(g *game.Data, unitId, eventPhase, expectedPhase string) bool {
	// Check if we're in the expected phase
	currentPhaseName := getPhaseString(g.CurrentPhase)
	if currentPhaseName != expectedPhase {
		return false
	}
	
	// Check if it's the right unit's turn
	if g.CurrentUnitIdx >= len(g.CurrentUnitOrder) {
		return false
	}
	
	currentUnitId := g.CurrentUnitOrder[g.CurrentUnitIdx]
	if currentUnitId != unitId {
		return false
	}
	
	// Verify the phase in the event matches expected phase
	return eventPhase == expectedPhase
}

// getPhaseString converts Phase enum to string
func getPhaseString(phase game.Phase) string {
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