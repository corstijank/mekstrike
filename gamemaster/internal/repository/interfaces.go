package repository

import (
	"context"

	"github.com/corstijank/mekstrike/gamemaster/game"
)

type GameRepository interface {
	Save(ctx context.Context, gameData *game.Data) error
	Get(ctx context.Context, gameID string) (*game.Data, error)
	GetAll(ctx context.Context) ([]*game.Data, error)
	Delete(ctx context.Context, gameID string) error
	Exists(ctx context.Context, gameID string) (bool, error)
}