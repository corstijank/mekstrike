package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/corstijank/mekstrike/src/gamemaster/clients/armybuilder"
	"github.com/corstijank/mekstrike/src/gamemaster/clients/battlefield"
	"github.com/corstijank/mekstrike/src/gamemaster/clients/unit"
	"github.com/corstijank/mekstrike/src/gamemaster/game"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace/propagation"
	"google.golang.org/grpc/metadata"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type NewGameRequest struct {
	PlayerName string
}

var games map[string]*game.Data
var abc armybuilder.AbClient

func main() {
	var err error
	abc, err = armybuilder.New()
	if err != nil {
		log.Fatal(err)
	}
	defer abc.Close()

	games = make(map[string]*game.Data)

	r := mux.NewRouter()

	r.HandleFunc("/games/{id}", getGame).Methods("GET")
	r.HandleFunc("/games/{id}/currentOpts", getCurrentOptions).Methods("GET")
	r.HandleFunc("/games", newGame).Methods("POST")
	r.HandleFunc("/games", getGames).Methods("GET")

	log.Printf("Starting Mekstrike gamemaster")
	log.Fatal(http.ListenAndServe(":7011", r))
}

func getGame(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	writeJSONResponse(games[id], rw)
}

func getGames(rw http.ResponseWriter, r *http.Request) {
	result := make([]game.Data, 0)
	for _, v := range games {
		result = append(result, *v)
	}
	writeJSONResponse(result, rw)
}

func getCurrentOptions(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	options, err := games[id].GetCurrentMoveOptions(r.Context())
	if err != nil {
		http.Error(rw, "Oh no", 500)
	}
	writeJSONResponse(options, rw)
}

func newGame(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := tracecontext.HTTPFormat{}
	sc, _ := f.SpanContextFromRequest(r)
	traceContextBinary := propagation.Binary(sc)
	ctx = metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", string(traceContextBinary))

	log.Printf("Gamemaster::newGame - called as part of trace %+v", sc.TraceID)

	req := NewGameRequest{}
	if r.Body == nil {
		http.Error(rw, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		log.Println(err)
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}

	armyA, err := abc.CreateArmy(ctx, 1, 2, 1, 0)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		log.Println(err)
		return
	}
	armyB, err := abc.CreateArmy(ctx, 1, 2, 1, 0)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		log.Println(err)
		return
	}
	log.Printf("Force A: %+v\n", armyA)
	log.Printf("Force B: %+v\n", armyB)

	bf, err := battlefield.GetBattlefieldClient(id.String())
	if err != nil {
		log.Println(err)
	}
	// Should make an init here
	_, err = bf.GetBoardCells(ctx)
	if err != nil {
		log.Println(err)
	}
	playerAUnits := make([]string, 0)
	valueA := 0
	for _, u := range armyA {
		valueA += int(u.Pointvalue)
		// Should make DaprClient a parameter to skip useless inits in client
		newUnit, _ := unit.NewUnit(bf.ID(), req.PlayerName, u)
		newUnit.Deploy(ctx, unit.DeployData{BattlefieldID: bf.ID(), Owner: req.PlayerName, Stats: u, DeployLocation: "NE"})
		playerAUnits = append(playerAUnits, newUnit.ID())
	}

	playerBUnits := make([]string, 0)
	valueB := 0
	for _, u := range armyB {
		valueB += int(u.Pointvalue)
		newUnit, _ := unit.NewUnit(bf.ID(), "CPU", u)
		newUnit.Deploy(ctx, unit.DeployData{BattlefieldID: bf.ID(), Owner: "CPU", Stats: u, DeployLocation: "SW"})
		playerBUnits = append(playerBUnits, newUnit.ID())
	}

	g := game.Data{
		ID:            id.String(),
		StartTime:     time.Now(),
		PlayerA:       req.PlayerName,
		PlayerAValue:  valueA,
		PlayerAUnits:  playerAUnits,
		PlayerB:       "CPU",
		PlayerBValue:  valueB,
		PlayerBUnits:  playerBUnits,
		Battlefieldld: bf.ID(),
		ActivePlayer:  game.PlayerA,
		CurrentRound:  0,
		CurrentPhase:  game.Movement,
	}
	games[g.ID] = &g
	g.StartGame(ctx)

	writeJSONResponse(g, rw)
}

func writeJSONResponse(obj interface{}, w http.ResponseWriter) {
	js, err := json.Marshal(obj)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
