package types

import (
	"encoding/json"
	"testing"
)

func TestNewGameRequestJSON(t *testing.T) {
	req := NewGameRequest{PlayerName: "Alice"}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if got, want := string(b), `{"playerName":"Alice"}`; got != want {
		t.Fatalf("json = %s, want %s", got, want)
	}

	var roundtrip NewGameRequest
	if err := json.Unmarshal(b, &roundtrip); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if roundtrip.PlayerName != "Alice" {
		t.Fatalf("PlayerName = %q, want %q", roundtrip.PlayerName, "Alice")
	}
}

func TestMoveRequestJSON(t *testing.T) {
	req := MoveRequest{UnitId: "U1", X: 3, Y: 4, Heading: 2}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if got, want := string(b), `{"unitId":"U1","x":3,"y":4,"heading":2}`; got != want {
		t.Fatalf("json = %s, want %s", got, want)
	}

	var rt MoveRequest
	if err := json.Unmarshal(b, &rt); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if rt.UnitId != "U1" || rt.X != 3 || rt.Y != 4 || rt.Heading != 2 {
		t.Fatalf("roundtrip mismatch: %+v", rt)
	}
}
