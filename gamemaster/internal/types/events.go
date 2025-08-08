package types

type DaprSubscription struct {
	PubsubName string `json:"pubsubname"`
	Topic      string `json:"topic"`
	Route      string `json:"route"`
}

type ActionCompletedEvent struct {
	GameId         string            `json:"GameId"`
	UnitId         string            `json:"UnitId"`
	Phase          string            `json:"Phase"`
	BattlefieldId  string            `json:"BattlefieldId,omitempty"`
	SourceLocation map[string]int    `json:"SourceLocation,omitempty"`
	TargetLocation map[string]int    `json:"TargetLocation,omitempty"`
	TargetId       string            `json:"TargetId,omitempty"`
	Unit           map[string]interface{} `json:"Unit,omitempty"`
}

type CloudEvent struct {
	Data            string `json:"data"`
	DataContentType string `json:"datacontenttype"`
	ID              string `json:"id"`
	PubsubName      string `json:"pubsubname"`
	Source          string `json:"source"`
	SpecVersion     string `json:"specversion"`
	Time            string `json:"time"`
	Topic           string `json:"topic"`
	Type            string `json:"type"`
}