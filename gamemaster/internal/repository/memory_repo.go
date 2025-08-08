package repository

import (
	"context"
	"sort"
	"sync"

	"github.com/corstijank/mekstrike/gamemaster/game"
	"github.com/corstijank/mekstrike/gamemaster/internal/types"
)

type MemoryGameRepository struct {
	games map[string]*game.Data
	mutex sync.RWMutex
}

func NewMemoryGameRepository() *MemoryGameRepository {
	return &MemoryGameRepository{
		games: make(map[string]*game.Data),
	}
}

func (r *MemoryGameRepository) Save(ctx context.Context, gameData *game.Data) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.games[gameData.ID] = gameData
	return nil
}

func (r *MemoryGameRepository) Get(ctx context.Context, gameID string) (*game.Data, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	gameData, exists := r.games[gameID]
	if !exists {
		return nil, types.NewGameNotFoundError(gameID)
	}
	
	return gameData, nil
}

func (r *MemoryGameRepository) GetAll(ctx context.Context) ([]*game.Data, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	result := make([]*game.Data, 0, len(r.games))
	for _, gameData := range r.games {
		result = append(result, gameData)
	}
	
	// Sort by start time
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartTime.Before(result[j].StartTime)
	})
	
	return result, nil
}

func (r *MemoryGameRepository) Delete(ctx context.Context, gameID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.games, gameID)
	return nil
}

func (r *MemoryGameRepository) Exists(ctx context.Context, gameID string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	_, exists := r.games[gameID]
	return exists, nil
}