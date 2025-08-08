package types

import "fmt"

type GameError struct {
	Code    int
	Message string
	Err     error
}

func (e *GameError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewGameNotFoundError(gameID string) *GameError {
	return &GameError{
		Code:    404,
		Message: fmt.Sprintf("Game not found: %s", gameID),
	}
}

func NewInvalidRequestError(message string, err error) *GameError {
	return &GameError{
		Code:    400,
		Message: message,
		Err:     err,
	}
}

func NewInternalError(message string, err error) *GameError {
	return &GameError{
		Code:    500,
		Message: message,
		Err:     err,
	}
}