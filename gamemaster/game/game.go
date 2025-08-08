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

type AvailableActions struct {
	Options
	UnitOwner          string
	ActionType         string
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

func (g Data) GetAvailableActions(ctx context.Context, client dapr.Client) (AvailableActions, error) {
	currentUnitID := g.CurrentUnitOrder[g.CurrentUnitIdx]
	currentUnit := uclient.GetUnitClient(client, currentUnitID)
	unitData, err := currentUnit.GetData(ctx)
	if err != nil {
		return AvailableActions{}, err
	}

	bf := bfclient.GetBattlefieldClient(client, g.Battlefieldld)
	options, err := bf.GetMovementOptions(ctx, &unitData)
	if err != nil {
		return AvailableActions{}, err
	}

	return AvailableActions{
		Options: Options{
			GameID:        g.ID,
			CurrentRound:  g.CurrentRound,
			CurrentPhase:  g.CurrentPhase,
			CurrentUnitID: currentUnitID,
		},
		UnitOwner:          unitData.Owner,
		ActionType:         "movement",
		AllowedCoordinates: options}, nil

}

func (g *Data) NewRound(ctx context.Context, client dapr.Client) {
	g.CurrentRound++
	g.CurrentPhase = Movement
	g.CurrentUnitIdx = 0
	log.Printf("New Round: %d", g.CurrentRound)
	g.activateCurrentUnit(ctx, client)
}

func (g *Data) activateCurrentUnit(ctx context.Context, client dapr.Client) {
	if g.CurrentUnitIdx >= len(g.CurrentUnitOrder) {
		log.Printf("No more units in order, advancing phase or round")
		return
	}

	currentUnitID := g.CurrentUnitOrder[g.CurrentUnitIdx]
	log.Printf("Activating unit %s", currentUnitID)
	
	u := uclient.GetUnitClient(client, currentUnitID)
	err := u.SetActive(ctx, true)
	if err != nil {
		log.Printf("Error setting active\n%v", err.Error())
		return
	}

	// Check if current unit is AI-owned and publish event
	if g.isAIUnit(currentUnitID) {
		err := g.publishAITurnEvent(ctx, client, currentUnitID)
		if err != nil {
			log.Printf("Error publishing AI turn event: %v", err)
		}
	}
}

func (g *Data) isAIUnit(unitID string) bool {
	// Check if unit belongs to CPU player
	for _, cpuUnit := range g.PlayerBUnits {
		if cpuUnit == unitID {
			return true
		}
	}
	return false
}

func (g *Data) publishAITurnEvent(ctx context.Context, client dapr.Client, unitID string) error {
	eventData := map[string]interface{}{
		"gameId": g.ID,
		"unitId": unitID,
		"phase":  g.phaseToString(),
		"round":  g.CurrentRound,
	}

	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	err = client.PublishEvent(ctx, "redis-pubsub", "ai-turn-started", eventJSON)
	if err != nil {
		return err
	}

	log.Printf("Published AI turn event: %s", string(eventJSON))
	return nil
}

func (g *Data) phaseToString() string {
	switch g.CurrentPhase {
	case Movement:
		return "Movement"
	case Combat:
		return "Combat"
	case End:
		return "End"
	default:
		return "Unknown"
	}
}

func (g *Data) AdvanceTurn(ctx context.Context, client dapr.Client) {
	// Deactivate current unit
	if g.CurrentUnitIdx < len(g.CurrentUnitOrder) {
		currentUnitID := g.CurrentUnitOrder[g.CurrentUnitIdx]
		u := uclient.GetUnitClient(client, currentUnitID)
		err := u.SetActive(ctx, false)
		if err != nil {
			log.Printf("Error deactivating unit %s: %v", currentUnitID, err)
		}
	}

	// Advance to next unit
	g.CurrentUnitIdx++
	
	// Check if all units have moved in current phase
	if g.CurrentUnitIdx >= len(g.CurrentUnitOrder) {
		g.advancePhase(ctx, client)
	} else {
		g.activateCurrentUnit(ctx, client)
	}
}

func (g *Data) advancePhase(ctx context.Context, client dapr.Client) {
	switch g.CurrentPhase {
	case Movement:
		log.Printf("Advancing to Combat phase")
		g.CurrentPhase = Combat
		g.CurrentUnitIdx = 0
		g.activateCurrentUnit(ctx, client)
	case Combat:
		log.Printf("Advancing to End phase")
		g.CurrentPhase = End
		g.CurrentUnitIdx = 0
		g.activateCurrentUnit(ctx, client)
	case End:
		log.Printf("Round completed, starting new round")
		g.NewRound(ctx, client)
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
