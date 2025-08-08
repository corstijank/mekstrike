package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/corstijank/mekstrike/gamemaster/clients/bfclient"
	"github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	"github.com/corstijank/mekstrike/gamemaster/game"
	dapr "github.com/dapr/go-sdk/client"

	"github.com/gorilla/mux"
)

type NewGameRequest struct {
	PlayerName string
}

var client dapr.Client
var games map[string]*game.Data

func main() {
	client, _ = dapr.NewClient()

	games = make(map[string]*game.Data)

	r := mux.NewRouter()

	r.HandleFunc("/games/{id}", getGame).Methods("GET")
	r.HandleFunc("/games/{id}/availableActions", getAvailableActions).Methods("GET")
	r.HandleFunc("/games/{id}/board", getBoard).Methods("GET")
	r.HandleFunc("/games/{id}/units/{uid}", getUnit).Methods("GET")
	r.HandleFunc("/games/{id}/advanceTurn", advanceTurn).Methods("POST")
	r.HandleFunc("/games", newGame).Methods("POST")
	r.HandleFunc("/games", getGames).Methods("GET")
	
	// Dapr pub/sub event handlers
	r.HandleFunc("/dapr/subscribe", getDaprSubscriptions).Methods("GET")
	r.HandleFunc("/unit-movement-completed", handleMovementCompleted).Methods("POST")
	r.HandleFunc("/unit-attack-completed", handleAttackCompleted).Methods("POST")
	r.HandleFunc("/unit-end-phase-completed", handleEndPhaseCompleted).Methods("POST")

	log.Printf("Starting Mekstrike gamemaster")
	log.Fatal(http.ListenAndServe(":7011", r))
}

func getGame(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	writeJSONResponse(games[id], rw)
}

func getGames(rw http.ResponseWriter, r *http.Request) {
	result := make(gameList, 0)
	for _, v := range games {
		result = append(result, *v)
	}
	sort.Sort(result)
	writeJSONResponse(result, rw)
}

func getAvailableActions(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	options, err := games[id].GetAvailableActions(r.Context(), client)
	if err != nil {
		http.Error(rw, "Oh no", 500)
	}
	writeJSONResponse(options, rw)
}

func getBoard(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := mux.Vars(r)["id"]
	g := games[id]
	bfc := bfclient.GetBattlefieldClient(client, g.Battlefieldld)
	cells, err := bfc.GetBoardCells(ctx)
	if err != nil {
		http.Error(rw, "Error getting cells", 500)
	}
	cols, err := bfc.GetNumberOfCols(ctx)
	if err != nil {
		http.Error(rw, "Error getting cols", 500)
	}
	rows, err := bfc.GetNumberOfRows(ctx)
	if err != nil {
		http.Error(rw, "Error getting rows", 500)
	}
	writeJSONResponse(map[string]any{"cells": cells, "rows": rows, "cols": cols}, rw)
}

func getUnit(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uid := mux.Vars(r)["uid"]
	uc := uclient.GetUnitClient(client, uid)

	unitData, err := uc.GetData(ctx)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Error getting unit client", 500)
	}
	writeJSONResponse(unitData, rw)
}

func newGame(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Printf("Gamemaster::newGame - called")

	req := NewGameRequest{}
	if r.Body == nil {
		http.Error(rw, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		log.Println(err)
		return
	}
	g, err := game.New(ctx, client, req.PlayerName)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		log.Println(err)
		return
	}
	games[g.ID] = &g
	g.StartGame(ctx, client)

	writeJSONResponse(g, rw)
}

func advanceTurn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	
	g := games[id]
	if g == nil {
		http.Error(rw, "Game not found", 404)
		return
	}
	
	g.AdvanceTurn(ctx, client)
	writeJSONResponse(g, rw)
}

type DaprSubscription struct {
	PubsubName string `json:"pubsubname"`
	Topic      string `json:"topic"`
	Route      string `json:"route"`
}

type ActionCompletedEvent struct {
	GameId string `json:"GameId"`
	UnitId string `json:"UnitId"`
	Phase  string `json:"Phase"`
}

type CloudEvent struct {
	Data            string `json:"data"`
	DataContentType string `json:"datacontenttype"`
	ID              string `json:"id"`
	PubsubName      string `json:"pubsubname"`
	Source          string `json:"source"`
	SpecVersion     string `json:"specversion"`
	Time            string `json:"time"`
	Topic           string `json:"topic"`
	Type            string `json:"type"`
}

func parseActionCompletedEvent(r *http.Request) (*ActionCompletedEvent, error) {
	// Try to parse as CloudEvent first
	var cloudEvent CloudEvent
	body := make([]byte, 0)
	if data, err := io.ReadAll(r.Body); err == nil {
		body = data
		r.Body = io.NopCloser(bytes.NewReader(body)) // Reset body for potential re-read
	}

	if err := json.Unmarshal(body, &cloudEvent); err == nil && cloudEvent.Data != "" {
		// Successfully parsed as CloudEvent, now parse the nested data
		log.Printf("Parsing CloudEvent: data=%s, topic=%s, type=%s", cloudEvent.Data, cloudEvent.Topic, cloudEvent.Type)
		
		var event ActionCompletedEvent
		if err := json.Unmarshal([]byte(cloudEvent.Data), &event); err != nil {
			return nil, fmt.Errorf("failed to parse CloudEvent data field: %v", err)
		}
		
		log.Printf("Parsed CloudEvent data: %+v", event)
		return &event, nil
	}

	// Fall back to direct parsing (backward compatibility)
	log.Printf("Trying direct parsing of body: %s", string(body))
	var event ActionCompletedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse as both CloudEvent and direct ActionCompletedEvent: %v", err)
	}
	
	log.Printf("Parsed direct event: %+v", event)
	return &event, nil
}

func getDaprSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions := []DaprSubscription{
		{
			PubsubName: "redis-pubsub",
			Topic:      "unit-movement-completed",
			Route:      "unit-movement-completed",
		},
		{
			PubsubName: "redis-pubsub",
			Topic:      "unit-attack-completed",
			Route:      "unit-attack-completed",
		},
		{
			PubsubName: "redis-pubsub",
			Topic:      "unit-end-phase-completed",
			Route:      "unit-end-phase-completed",
		},
	}
	writeJSONResponse(subscriptions, w)
}

func handleMovementCompleted(w http.ResponseWriter, r *http.Request) {
	event, err := parseActionCompletedEvent(r)
	if err != nil {
		log.Printf("Error parsing movement completed event: %v", err)
		http.Error(w, "Invalid event data", 400)
		return
	}
	
	log.Printf("Received movement completed event: %+v", *event)
	
	// Idempotency check: verify this event matches current game state
	if g := games[event.GameId]; g != nil {
		if !isValidPhaseEvent(g, event.UnitId, event.Phase, "Movement") {
			log.Printf("Ignoring duplicate or out-of-order movement event for game %s, unit %s", event.GameId, event.UnitId)
			w.WriteHeader(200) // Still return success to avoid retries
			return
		}
		
		g.AdvanceTurn(r.Context(), client)
		log.Printf("Advanced turn for game %s after movement completion", event.GameId)
	} else {
		log.Printf("No game found with ID %s for movement completion", event.GameId)
	}
	
	w.WriteHeader(200)
}

func handleAttackCompleted(w http.ResponseWriter, r *http.Request) {
	event, err := parseActionCompletedEvent(r)
	if err != nil {
		log.Printf("Error parsing attack completed event: %v", err)
		http.Error(w, "Invalid event data", 400)
		return
	}
	
	log.Printf("Received attack completed event: %+v", *event)
	
	// Idempotency check: verify this event matches current game state
	if g := games[event.GameId]; g != nil {
		if !isValidPhaseEvent(g, event.UnitId, event.Phase, "Combat") {
			log.Printf("Ignoring duplicate or out-of-order attack event for game %s, unit %s", event.GameId, event.UnitId)
			w.WriteHeader(200) // Still return success to avoid retries
			return
		}
		
		g.AdvanceTurn(r.Context(), client)
		log.Printf("Advanced turn for game %s after attack completion", event.GameId)
	} else {
		log.Printf("No game found with ID %s for attack completion", event.GameId)
	}
	
	w.WriteHeader(200)
}

func handleEndPhaseCompleted(w http.ResponseWriter, r *http.Request) {
	event, err := parseActionCompletedEvent(r)
	if err != nil {
		log.Printf("Error parsing end phase completed event: %v", err)
		http.Error(w, "Invalid event data", 400)
		return
	}
	
	log.Printf("Received end phase completed event: %+v", *event)
	
	// Idempotency check: verify this event matches current game state
	if g := games[event.GameId]; g != nil {
		if !isValidPhaseEvent(g, event.UnitId, event.Phase, "End") {
			log.Printf("Ignoring duplicate or out-of-order end phase event for game %s, unit %s", event.GameId, event.UnitId)
			w.WriteHeader(200) // Still return success to avoid retries
			return
		}
		
		g.AdvanceTurn(r.Context(), client)
		log.Printf("Advanced turn for game %s after end phase completion", event.GameId)
	} else {
		log.Printf("No game found with ID %s for end phase completion", event.GameId)
	}
	
	w.WriteHeader(200)
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

func writeJSONResponse(obj any, w http.ResponseWriter) {
	js, err := json.Marshal(obj)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type gameList []game.Data

func (gl gameList) Len() int {
	return len(gl)
}

func (gl gameList) Less(i, j int) bool {
	return gl[i].StartTime.Before(gl[j].StartTime)
}

func (gl gameList) Swap(i, j int) {
	gl[i], gl[j] = gl[j], gl[i]
}
