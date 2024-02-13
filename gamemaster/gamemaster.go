package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/corstijank/mekstrike/gamemaster/clients/bfclient"
	"github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	"github.com/corstijank/mekstrike/gamemaster/game"
	dapr "github.com/dapr/go-sdk/client"

	"github.com/gorilla/mux"
)

type NewGameRequest struct {
	PlayerName string
}

var client dapr.Client
var games map[string]*game.Data

func main() {
	client, _ = dapr.NewClient()

	games = make(map[string]*game.Data)

	r := mux.NewRouter()

	r.HandleFunc("/games/{id}", getGame).Methods("GET")
	r.HandleFunc("/games/{id}/currentOpts", getCurrentOptions).Methods("GET")
	r.HandleFunc("/games/{id}/board", getBoard).Methods("GET")
	r.HandleFunc("/games/{id}/units/{uid}", getUnit).Methods("GET")
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
	result := make(gameList, 0)
	for _, v := range games {
		result = append(result, *v)
	}
	sort.Sort(result)
	writeJSONResponse(result, rw)
}

func getCurrentOptions(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	options, err := games[id].GetCurrentMoveOptions(r.Context(), client)
	if err != nil {
		http.Error(rw, "Oh no", 500)
	}
	writeJSONResponse(options, rw)
}

func getBoard(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := mux.Vars(r)["id"]
	g := games[id]
	bfc := bfclient.GetBattlefieldClient(client, g.Battlefieldld)
	cells, err := bfc.GetBoardCells(ctx)
	if err != nil {
		http.Error(rw, "Error getting cells", 500)
	}
	cols, err := bfc.GetNumberOfCols(ctx)
	if err != nil {
		http.Error(rw, "Error getting cols", 500)
	}
	rows, err := bfc.GetNumberOfRows(ctx)
	if err != nil {
		http.Error(rw, "Error getting rows", 500)
	}
	writeJSONResponse(map[string]interface{}{"cells": cells, "rows": rows, "cols": cols}, rw)
}

func getUnit(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uid := mux.Vars(r)["uid"]
	uc := uclient.GetUnitClient(client, uid)

	unitData, err := uc.GetData(ctx)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Error getting unit client", 500)
	}
	writeJSONResponse(unitData, rw)
}

func newGame(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Printf("Gamemaster::newGame - called")

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
	g, err := game.New(ctx, client, req.PlayerName)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		log.Println(err)
		return
	}
	games[g.ID] = &g
	g.StartGame(ctx, client)

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

type gameList []game.Data

func (gl gameList) Len() int {
	return len(gl)
}

func (gl gameList) Less(i, j int) bool {
	return gl[i].StartTime.Before(gl[j].StartTime)
}

func (gl gameList) Swap(i, j int) {
	gl[i], gl[j] = gl[j], gl[i]
}
