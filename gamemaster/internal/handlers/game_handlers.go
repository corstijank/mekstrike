package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/corstijank/mekstrike/gamemaster/clients/bfclient"
	"github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	"github.com/corstijank/mekstrike/gamemaster/game"
	"github.com/corstijank/mekstrike/gamemaster/internal/repository"
	"github.com/corstijank/mekstrike/gamemaster/internal/types"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
)

type GameHandlers struct {
	client     dapr.Client
	repository repository.GameRepository
}

func NewGameHandlers(client dapr.Client, repo repository.GameRepository) *GameHandlers {
	return &GameHandlers{
		client:     client,
		repository: repo,
	}
}

func (h *GameHandlers) GetGame(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	
	gameData, err := h.repository.Get(r.Context(), id)
	if err != nil {
		if gameErr, ok := err.(*types.GameError); ok {
			http.Error(rw, gameErr.Message, gameErr.Code)
			return
		}
		http.Error(rw, "Internal server error", 500)
		return
	}
	
	writeJSONResponse(gameData, rw)
}

func (h *GameHandlers) GetGames(rw http.ResponseWriter, r *http.Request) {
	games, err := h.repository.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting all games: %v", err)
		http.Error(rw, "Internal server error", 500)
		return
	}
	
	// Convert to slice for JSON response
	result := make([]game.Data, 0, len(games))
	for _, gameData := range games {
		result = append(result, *gameData)
	}
	
	writeJSONResponse(result, rw)
}

func (h *GameHandlers) GetAvailableActions(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	
	gameData, err := h.repository.Get(r.Context(), id)
	if err != nil {
		if gameErr, ok := err.(*types.GameError); ok {
			http.Error(rw, gameErr.Message, gameErr.Code)
			return
		}
		http.Error(rw, "Internal server error", 500)
		return
	}
	
	options, err := gameData.GetAvailableActions(r.Context(), h.client)
	if err != nil {
		log.Printf("Error getting available actions: %v", err)
		http.Error(rw, "Error getting available actions", 500)
		return
	}
	
	writeJSONResponse(options, rw)
}

func (h *GameHandlers) GetBoard(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	
	gameData, err := h.repository.Get(ctx, id)
	if err != nil {
		if gameErr, ok := err.(*types.GameError); ok {
			http.Error(rw, gameErr.Message, gameErr.Code)
			return
		}
		http.Error(rw, "Internal server error", 500)
		return
	}
	
	bfc := bfclient.GetBattlefieldClient(h.client, gameData.Battlefieldld)
	
	cells, err := bfc.GetBoardCells(ctx)
	if err != nil {
		log.Printf("Error getting board cells: %v", err)
		http.Error(rw, "Error getting board cells", 500)
		return
	}
	
	cols, err := bfc.GetNumberOfCols(ctx)
	if err != nil {
		log.Printf("Error getting board cols: %v", err)
		http.Error(rw, "Error getting board dimensions", 500)
		return
	}
	
	rows, err := bfc.GetNumberOfRows(ctx)
	if err != nil {
		log.Printf("Error getting board rows: %v", err)
		http.Error(rw, "Error getting board dimensions", 500)
		return
	}
	
	writeJSONResponse(map[string]any{
		"cells": cells,
		"rows":  rows,
		"cols":  cols,
	}, rw)
}

func (h *GameHandlers) GetUnit(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := mux.Vars(r)["uid"]
	
	uc := uclient.GetUnitClient(h.client, uid)
	unitData, err := uc.GetData(ctx)
	if err != nil {
		log.Printf("Error getting unit data: %v", err)
		http.Error(rw, "Error getting unit data", 500)
		return
	}
	
	writeJSONResponse(unitData, rw)
}

func (h *GameHandlers) NewGame(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	log.Printf("Gamemaster::newGame - called")
	
	req := types.NewGameRequest{}
	if r.Body == nil {
		http.Error(rw, "Please send a request body", 400)
		return
	}
	
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(rw, "Invalid request body", 400)
		return
	}
	
	gameData, err := game.New(ctx, h.client, req.PlayerName)
	if err != nil {
		log.Printf("Error creating game: %v", err)
		http.Error(rw, "Error creating game", 500)
		return
	}
	
	err = h.repository.Save(ctx, &gameData)
	if err != nil {
		log.Printf("Error saving game: %v", err)
		http.Error(rw, "Error saving game", 500)
		return
	}
	
	gameData.StartGame(ctx, h.client)
	
	// Save updated game state after starting
	err = h.repository.Save(ctx, &gameData)
	if err != nil {
		log.Printf("Error saving started game: %v", err)
	}
	
	writeJSONResponse(gameData, rw)
}

func (h *GameHandlers) AdvanceTurn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	
	gameData, err := h.repository.Get(ctx, id)
	if err != nil {
		if gameErr, ok := err.(*types.GameError); ok {
			http.Error(rw, gameErr.Message, gameErr.Code)
			return
		}
		http.Error(rw, "Internal server error", 500)
		return
	}
	
	gameData.AdvanceTurn(ctx, h.client)
	
	err = h.repository.Save(ctx, gameData)
	if err != nil {
		log.Printf("Error saving game after turn advance: %v", err)
	}
	
	writeJSONResponse(gameData, rw)
}

func (h *GameHandlers) GetGameLogs(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	
	gameData, err := h.repository.Get(r.Context(), id)
	if err != nil {
		if gameErr, ok := err.(*types.GameError); ok {
			http.Error(rw, gameErr.Message, gameErr.Code)
			return
		}
		http.Error(rw, "Internal server error", 500)
		return
	}
	
	// Transform CloudEvents to UI-compatible format
	messages := make([]map[string]interface{}, 0, len(gameData.GameLogs))
	for _, event := range gameData.GameLogs {
		messages = append(messages, map[string]interface{}{
			"type":      event.Type,
			"message":   event.Data,
			"timestamp": event.Time,
		})
	}
	
	writeJSONResponse(messages, rw)
}

func writeJSONResponse(obj any, w http.ResponseWriter) {
	js, err := json.Marshal(obj)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	
	_, err = w.Write(js)
	if err != nil {
		log.Printf("Write error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}