package game

import (
	"time"

	"github.com/corstijank/mekstrike/domain/battlefield"
)

type Player int

const (
	PlayerA Player = iota
	PlayerB
)

type Phase int

const (
	Movement Phase = iota
	Combat
	End
)

type Data struct {
	ID               string
	StartTime        time.Time
	PlayerA          string
	PlayerAUnits     []string
	PlayerAValue     int
	PlayerB          string
	PlayerBValue     int
	PlayerBUnits     []string
	Battlefieldld    string
	ActivePlayer     Player
	CurrentRound     int
	CurrentPhase     Phase
	CurrentUnitIdx   int
	CurrentUnitOrder []string
}

type Options struct {
	GameID        string
	CurrentRound  int
	CurrentPhase  Phase
	CurrentUnitID string
}

type AvailableActions struct {
	Options
	UnitOwner          string
	ActionType         string
	AllowedCoordinates []battlefield.Coordinates
}

type NewGameRequest struct {
	PlayerName string
}