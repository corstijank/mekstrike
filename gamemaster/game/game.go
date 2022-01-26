package game

import (
	"encoding/json"
	"time"
)

type Data struct {
	ID            string
	StartTime     time.Time
	PlayerA       string
	PlayerAUnits  []string
	PlayerAValue  int
	PlayerB       string
	PlayerBValue  int
	PlayerBUnits  []string
	BattlefieldID string
}

func (d Data) Marshal() ([]byte, error) {
	result, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d Data) GetKey() string {
	return d.ID
}

func (d Data) GetIndices() []string {
	return []string{"_games"}
}
