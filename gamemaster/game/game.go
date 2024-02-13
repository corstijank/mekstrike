package game

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/corstijank/mekstrike/domain/battlefield"
	unitType "github.com/corstijank/mekstrike/domain/unit"
	"github.com/corstijank/mekstrike/gamemaster/clients/armybuilder"
	"github.com/corstijank/mekstrike/gamemaster/clients/bfclient"
	"github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	dapr "github.com/dapr/go-sdk/client"
	uuid "github.com/satori/go.uuid"
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
	AllowedCoordinates []battlefield.Coordinates
}

type NewGameRequest struct {
	PlayerName string
}

func New(ctx context.Context, client dapr.Client, playername string) (Data, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}

	abc, err := armybuilder.New(client)
	if err != nil {
		return Data{}, err
	}
	armyA, err := abc.CreateArmy(ctx, 1, 2, 1, 0)
	if err != nil {
		return Data{}, err
	}
	armyB, err := abc.CreateArmy(ctx, 1, 2, 1, 0)
	if err != nil {
		return Data{}, err
	}

	log.Printf("Force A: %+v\n", armyA)
	log.Printf("Force B: %+v\n", armyB)

	bf := bfclient.GetBattlefieldClient(client, id.String())
	if err != nil {
		return Data{}, err
	}
	// Should make an init here
	_, err = bf.GetBoardCells(ctx)
	if err != nil {
		return Data{}, err
	}
	playerAUnits := make([]string, 0)
	valueA := 0
	for _, u := range armyA {
		valueA += int(u.Pointvalue)
		// Should make DaprClient a parameter to skip useless inits in client
		newUnit := uclient.NewUnit(client, bf.ID(), playername, u)
		err = newUnit.Deploy(ctx, unitType.DeployRequest{BattlefieldId: bf.ID(), Owner: playername, Stats: u, Corner: "NE"})
		if err != nil {
			return Data{}, err
		}
		playerAUnits = append(playerAUnits, newUnit.ID())
	}

	playerBUnits := make([]string, 0)
	valueB := 0
	for _, u := range armyB {
		valueB += int(u.Pointvalue)
		newUnit := uclient.NewUnit(client, bf.ID(), "CPU", u)
		err = newUnit.Deploy(ctx, unitType.DeployRequest{BattlefieldId: bf.ID(), Owner: "CPU", Stats: u, Corner: "SW"})
		if err != nil {
			return Data{}, err
		}
		playerBUnits = append(playerBUnits, newUnit.ID())
	}

	return Data{
		ID:            id.String(),
		StartTime:     time.Now(),
		PlayerA:       playername,
		PlayerAValue:  valueA,
		PlayerAUnits:  playerAUnits,
		PlayerB:       "CPU",
		PlayerBValue:  valueB,
		PlayerBUnits:  playerBUnits,
		Battlefieldld: bf.ID(),
		ActivePlayer:  PlayerA,
		CurrentRound:  0,
		CurrentPhase:  Movement,
	}, nil
}

func (g *Data) StartGame(ctx context.Context, client dapr.Client) {
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

	g.NewRound(ctx, client)

}

func (g Data) GetCurrentMoveOptions(ctx context.Context, client dapr.Client) (MoveOptions, error) {
	currentUnitID := g.CurrentUnitOrder[g.CurrentUnitIdx]
	currentUnit := uclient.GetUnitClient(client, currentUnitID)
	unitData, err := currentUnit.GetData(ctx)
	if err != nil {
		return MoveOptions{}, err
	}

	bf := bfclient.GetBattlefieldClient(client, g.Battlefieldld)
	options, err := bf.GetMovementOptions(ctx, &unitData)
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
		AllowedCoordinates: options}, nil

}

func (g *Data) NewRound(ctx context.Context, client dapr.Client) {
	g.CurrentRound++
	g.CurrentPhase = Movement
	g.CurrentUnitIdx = 0
	log.Printf("New Round: %d", g.CurrentRound)
	log.Printf("Activating unit %s", g.CurrentUnitOrder[g.CurrentUnitIdx])
	u := uclient.GetUnitClient(client, g.CurrentUnitOrder[g.CurrentUnitIdx])
	err := u.SetActive(ctx, true)
	if err != nil {
		log.Printf("Error setting active\n%v", err.Error())
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
