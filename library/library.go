package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/corstijank/mekstrike/src/common/go/storage"
	"github.com/corstijank/mekstrike/src/common/go/unit"
	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace/propagation"
	"google.golang.org/grpc/metadata"
)

var store storage.Store

func main() {
	log.Printf("### Starting Mekstrike Library")

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	var err error
	store, err = storage.New("library-store")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	r := mux.NewRouter()

	r.HandleFunc("/units", units).Methods("GET")
	r.HandleFunc("/units/by/{type}", unitsByType).Methods("GET")
	r.HandleFunc("/units/by/{type}/{class}", unitsByTypeAndClass).Methods("GET")
	r.HandleFunc("/units/by/{type}/{class}/random", randomUnitsFromTypeAndClass).Methods("GET")
	log.Fatal(http.ListenAndServe(":7010", r))
}

func units(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := tracecontext.HTTPFormat{}
	sc, _ := f.SpanContextFromRequest(r)
	traceContextBinary := propagation.Binary(sc)
	ctx = metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", string(traceContextBinary))

	log.Printf("Library::units - called as part of trace %+v", sc.TraceID)

	ir, err := store.ReadMany(ctx, "_units", &unit.UnitStats{})
	if err != nil {
		log.Printf("Error reading from store: %+v", err)
	}
	result, err := toUnitStats(ir)
	if err != nil {
		log.Printf("Error translating to unitstats: %+v", err)
	}
	writeJSONResponse(&result, rw)
}

func unitsByType(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := tracecontext.HTTPFormat{}
	sc, _ := f.SpanContextFromRequest(r)
	traceContextBinary := propagation.Binary(sc)
	ctx = metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", string(traceContextBinary))

	log.Printf("Library::unitsByType - called as part of trace %+v", sc.TraceID)

	t := mux.Vars(r)["type"]

	if t != "BM" {
		http.Error(rw, fmt.Sprintf("Unsupported unit type: %s", t), 500)
	}

	ir, err := store.ReadMany(ctx, "_units"+t, &unit.UnitStats{})
	if err != nil {
		log.Printf("Error reading from store: %+v", err)
	}
	result, err := toUnitStats(ir)
	if err != nil {
		log.Printf("Error translating to unitstats: %+v", err)
	}
	writeJSONResponse(&result, rw)
}

func randomUnitsFromTypeAndClass(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := tracecontext.HTTPFormat{}
	sc, _ := f.SpanContextFromRequest(r)
	traceContextBinary := propagation.Binary(sc)
	ctx = metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", string(traceContextBinary))

	log.Printf("Library:randomUnitsFromTypeAndClass - called as part of trace %+v", sc.TraceID)

	t := mux.Vars(r)["type"]
	c := mux.Vars(r)["class"]
	if t != "BM" {
		http.Error(rw, fmt.Sprintf("Unsupported unit type: %s", t), 500)
		return
	}
	s := 0
	if strings.ToLower(c) == "light" {
		s = 1
	} else if strings.ToLower(c) == "medium" {
		s = 2
	} else if strings.ToLower(c) == "heavy" {
		s = 3
	} else if strings.ToLower(c) == "assault" {
		s = 4
	} else {
		http.Error(rw, fmt.Sprintf("Unsupported unit class: %s", c), 500)
		return
	}
	ir, err := store.ReadMany(ctx, fmt.Sprintf("_units_%s_%d", t, s), &unit.UnitStats{})
	if err != nil {
		log.Printf("Error reading from store: %+v", err)
	}
	result, err := toUnitStats(ir)
	if err != nil {
		log.Printf("Error translating to unitstats: %+v", err)
	}
	writeJSONResponse(&result[rand.Intn(len(result))], rw)
}

func unitsByTypeAndClass(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := tracecontext.HTTPFormat{}
	sc, _ := f.SpanContextFromRequest(r)
	traceContextBinary := propagation.Binary(sc)
	ctx = metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", string(traceContextBinary))

	log.Printf("Library::unitsByTypeAndClass - called as part of trace %+v", sc.TraceID)

	t := mux.Vars(r)["type"]
	c := mux.Vars(r)["class"]
	if t != "BM" {
		http.Error(rw, fmt.Sprintf("Unsupported unit type: %s", t), 500)
		return
	}
	s := 0
	if strings.ToLower(c) == "light" {
		s = 1
	} else if strings.ToLower(c) == "medium" {
		s = 2
	} else if strings.ToLower(c) == "heavy" {
		s = 3
	} else if strings.ToLower(c) == "assault" {
		s = 4
	} else {
		http.Error(rw, fmt.Sprintf("Unsupported unit class: %s", c), 500)
		return
	}

	ir, err := store.ReadMany(ctx, fmt.Sprintf("_units_%s_%d", t, s), &unit.UnitStats{})
	if err != nil {
		log.Printf("Error reading from store: %+v", err)
	}
	result, err := toUnitStats(ir)
	if err != nil {
		log.Printf("Error translating to unitstats: %+v", err)
	}
	writeJSONResponse(&result, rw)
}

func toUnitStats(in []storage.Readable) ([]*unit.UnitStats, error) {
	result := make([]*unit.UnitStats, len(in))
	for i, o := range in {
		obj, ok := o.(*unit.UnitStats)
		if !ok {
			return nil, fmt.Errorf("%+v is not of type unit.Stats", o)
		}
		result[i] = obj
	}
	return result, nil
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
