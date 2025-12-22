package game

import "testing"

func TestIsAIUnit(t *testing.T) {
	g := &Data{PlayerBUnits: []string{"U1", "U2"}}
	if !g.isAIUnit("U1") {
		t.Fatalf("expected AI unit to be true for U1")
	}
	if g.isAIUnit("X") {
		t.Fatalf("expected AI unit to be false for X")
	}
}

func TestPhaseToString(t *testing.T) {
	g := &Data{}
	g.CurrentPhase = Movement
	if got := g.phaseToString(); got != "Movement" {
		t.Fatalf("got %q, want %q", got, "Movement")
	}
	g.CurrentPhase = Combat
	if got := g.phaseToString(); got != "Combat" {
		t.Fatalf("got %q, want %q", got, "Combat")
	}
	g.CurrentPhase = End
	if got := g.phaseToString(); got != "End" {
		t.Fatalf("got %q, want %q", got, "End")
	}
	g.CurrentPhase = Phase(99)
	if got := g.phaseToString(); got != "Unknown" {
		t.Fatalf("got %q, want %q", got, "Unknown")
	}
}
