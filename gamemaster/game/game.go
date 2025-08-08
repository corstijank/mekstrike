package game

import (
	"encoding/json"
	"time"

	"github.com/corstijank/mekstrike/gamemaster/internal/types"
	uuid "github.com/satori/go.uuid"
)

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

func (d *Data) AddGameEvent(eventType, source, data string) {
	id, _ := uuid.NewV4()
	event := types.CloudEvent{
		Data:            data,
		DataContentType: "application/json",
		ID:              id.String(),
		Source:          source,
		SpecVersion:     "1.0",
		Time:            time.Now().Format(time.RFC3339),
		Type:            eventType,
	}
	d.GameLogs = append(d.GameLogs, event)
}

func (d *Data) LogSystemEvent(message string) {
	d.AddGameEvent("system", "gamemaster", message)
}

func (d *Data) LogMovementEvent(message string) {
	d.AddGameEvent("movement", "gamemaster", message)
}

func (d *Data) LogCombatEvent(message string) {
	d.AddGameEvent("combat", "gamemaster", message)
}

func (d *Data) LogErrorEvent(message string) {
	d.AddGameEvent("error", "gamemaster", message)
}