package game

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/corstijank/mekstrike/src/gamemaster/clients/battlefield"
	"github.com/corstijank/mekstrike/src/gamemaster/clients/unit"
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

type MoveOptions struct {
	Options
	AllowedCells []battlefield.Cell
}

func (g *Data) StartGame(ctx context.Context) {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	g.CurrentRound = 0

	unitsA := make([]string, 0)
	unitsA = append(unitsA, g.PlayerAUnits...)
	randomizer.Shuffle(len(unitsA), func(i, j int) { unitsA[i], unitsA[j] = unitsA[j], unitsA[i] })

	unitsB := make([]string, 0)
	unitsB = append(unitsB, g.PlayerBUnits...)
	randomizer.Shuffle(len(unitsB), func(i, j int) { unitsB[i], unitsB[j] = unitsB[j], unitsB[i] })

	// Roll for initiative!
	if randomizer.Intn(100) <= 50 {
		// Player A starts
		g.ActivePlayer = PlayerA

		unitOrder := make([]string, 0)
		for {
			var unitA, unitB string

			if len(unitsA) >= 1 {
				unitA, unitsA = unitsA[0], unitsA[1:]
				unitOrder = append(unitOrder, unitA)
			}
			if len(unitsB) >= 1 {
				unitB, unitsB = unitsB[0], unitsB[1:]
				unitOrder = append(unitOrder, unitB)
			}
			if len(unitsA) == 0 && len(unitsB) == 0 {
				break
			}
		}
		g.CurrentUnitOrder = unitOrder
	} else {
		//Player B starts
		g.ActivePlayer = PlayerB

		unitOrder := make([]string, 0)
		for {
			var unitA, unitB string

			if len(unitsB) >= 1 {
				unitB, unitsB = unitsB[0], unitsB[1:]
				unitOrder = append(unitOrder, unitB)
			}
			if len(unitsA) >= 1 {
				unitA, unitsA = unitsA[0], unitsA[1:]
				unitOrder = append(unitOrder, unitA)
			}
			if len(unitsA) == 0 && len(unitsB) == 0 {
				break
			}
		}
		g.CurrentUnitOrder = unitOrder
	}

	g.NewRound(ctx)

}

func (g Data) GetCurrentMoveOptions(ctx context.Context) (MoveOptions, error) {
	currentUnitID := g.CurrentUnitOrder[g.CurrentUnitIdx]
	currentUnit, err := unit.GetUnitClient(currentUnitID)
	if err != nil {
		return MoveOptions{}, err
	}
	unitData, err := currentUnit.GetData(ctx)
	if err != nil {
		return MoveOptions{}, err
	}

	bf, err := battlefield.GetBattlefieldClient(g.Battlefieldld)
	if err != nil {
		return MoveOptions{}, err
	}

	options, err := bf.GetMovementOptions(ctx, unitData)
	if err != nil {
		return MoveOptions{}, err
	}

	return MoveOptions{
		Options: Options{
			GameID:        g.ID,
			CurrentRound:  g.CurrentRound,
			CurrentPhase:  g.CurrentPhase,
			CurrentUnitID: currentUnitID,
		},
		AllowedCells: options}, nil

}

func (g *Data) NewRound(ctx context.Context) {
	g.CurrentRound++
	g.CurrentPhase = Movement
	g.CurrentUnitIdx = 0
	log.Printf("New Round: %d", g.CurrentRound)
	log.Printf("Activating unit %s", g.CurrentUnitOrder[g.CurrentUnitIdx])
	u, err := unit.GetUnitClient(g.CurrentUnitOrder[g.CurrentUnitIdx])
	if err != nil {
		log.Println(err.Error())
	}
	err = u.SetActive(ctx, true)
	if err != nil {
		log.Println(err.Error())
	}
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
