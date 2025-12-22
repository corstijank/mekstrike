package game

import (
	"context"
	"testing"
)

// These tests focus on state transitions, not external side effects.
func TestAdvanceTurn_AdvancesPhaseWhenEndOfOrder(t *testing.T) {
	g := &Data{
		CurrentPhase:     Movement,
		CurrentUnitIdx:   0,
		CurrentUnitOrder: []string{"U1"},
	}

	// client is not used for assertions; nil is acceptable for state-only checks
	g.AdvanceTurn(context.Background(), nil)

	if g.CurrentPhase != Combat {
		t.Fatalf("CurrentPhase = %v, want Combat", g.CurrentPhase)
	}
	if g.CurrentUnitIdx != 0 {
		t.Fatalf("CurrentUnitIdx = %d, want 0", g.CurrentUnitIdx)
	}
}

func TestAdvanceTurn_RoundProgression(t *testing.T) {
	g := &Data{
		CurrentPhase:     Movement,
		CurrentUnitIdx:   0,
		CurrentUnitOrder: []string{"U1"},
	}

	// Move from Movement -> Combat
	g.AdvanceTurn(context.Background(), nil)
	if g.CurrentPhase != Combat {
		t.Fatalf("phase = %v, want Combat", g.CurrentPhase)
	}

	// Move from Combat -> End
	g.AdvanceTurn(context.Background(), nil)
	if g.CurrentPhase != End {
		t.Fatalf("phase = %v, want End", g.CurrentPhase)
	}

	// Move from End -> New Round (Movement) and round increments
	g.AdvanceTurn(context.Background(), nil)
	if g.CurrentPhase != Movement {
		t.Fatalf("phase = %v, want Movement", g.CurrentPhase)
	}
	if g.CurrentRound != 1 {
		t.Fatalf("round = %d, want 1", g.CurrentRound)
	}
}
