package services

import (
	"context"
	"log"

	"github.com/corstijank/mekstrike/gamemaster/game"
	"github.com/corstijank/mekstrike/gamemaster/internal/repository"
	dapr "github.com/dapr/go-sdk/client"
)

type GameService struct {
	client     dapr.Client
	repository repository.GameRepository
}

func NewGameService(client dapr.Client, repo repository.GameRepository) *GameService {
	return &GameService{
		client:     client,
		repository: repo,
	}
}

func (s *GameService) CreateGame(ctx context.Context, playerName string) (*game.Data, error) {
	gameData, err := game.New(ctx, s.client, playerName)
	if err != nil {
		return nil, err
	}
	
	err = s.repository.Save(ctx, &gameData)
	if err != nil {
		return nil, err
	}
	
	gameData.StartGame(ctx, s.client)
	
	// Save updated game state after starting
	err = s.repository.Save(ctx, &gameData)
	if err != nil {
		log.Printf("Error saving started game: %v", err)
	}
	
	return &gameData, nil
}

func (s *GameService) GetGame(ctx context.Context, gameID string) (*game.Data, error) {
	return s.repository.Get(ctx, gameID)
}

func (s *GameService) GetAllGames(ctx context.Context) ([]*game.Data, error) {
	return s.repository.GetAll(ctx)
}

func (s *GameService) AdvanceGameTurn(ctx context.Context, gameID string) (*game.Data, error) {
	gameData, err := s.repository.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}
	
	gameData.AdvanceTurn(ctx, s.client)
	
	err = s.repository.Save(ctx, gameData)
	if err != nil {
		log.Printf("Error saving game after turn advance: %v", err)
	}
	
	return gameData, nil
}

func (s *GameService) GetAvailableActions(ctx context.Context, gameID string) (*game.AvailableActions, error) {
	gameData, err := s.repository.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}
	
	actions, err := gameData.GetAvailableActions(ctx, s.client)
	if err != nil {
		return nil, err
	}
	
	return &actions, nil
}