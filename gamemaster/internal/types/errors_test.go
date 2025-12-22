package types

import (
	"errors"
	"testing"
)

func TestGameErrorFormatting(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		err := NewGameNotFoundError("game-123")
		if got, want := err.Error(), "Game not found: game-123"; got != want {
			t.Fatalf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("with wrapped error", func(t *testing.T) {
		base := errors.New("boom")
		err := NewInvalidRequestError("invalid payload", base)
		if got, want := err.Error(), "invalid payload: boom"; got != want {
			t.Fatalf("Error() = %q, want %q", got, want)
		}
	})
}

func TestErrorConstructors(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		e := NewGameNotFoundError("g1")
		if e.Code != 404 {
			t.Errorf("Code = %d, want 404", e.Code)
		}
		if e.Message == "" {
			t.Errorf("Message should not be empty")
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		base := errors.New("details")
		e := NewInvalidRequestError("bad", base)
		if e.Code != 400 {
			t.Errorf("Code = %d, want 400", e.Code)
		}
		if e.Err == nil {
			t.Errorf("Err should be set")
		}
	})

	t.Run("internal error", func(t *testing.T) {
		base := errors.New("oops")
		e := NewInternalError("failed", base)
		if e.Code != 500 {
			t.Errorf("Code = %d, want 500", e.Code)
		}
		if e.Err == nil {
			t.Errorf("Err should be set")
		}
	})
}
