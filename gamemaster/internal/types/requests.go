package types

type NewGameRequest struct {
	PlayerName string `json:"playerName"`
}

type MoveRequest struct {
	UnitId  string `json:"unitId"`
	X       int32  `json:"x"`
	Y       int32  `json:"y"`
	Heading int32  `json:"heading"`
}