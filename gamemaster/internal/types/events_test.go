package types

import (
	"encoding/json"
	"testing"
)

func TestDaprSubscriptionJSON(t *testing.T) {
	sub := DaprSubscription{PubsubName: "redis-pubsub", Topic: "unit-movement-completed", Route: "/events/move"}
	b, err := json.Marshal(sub)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var rt DaprSubscription
	if err := json.Unmarshal(b, &rt); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if rt.PubsubName != sub.PubsubName || rt.Topic != sub.Topic || rt.Route != sub.Route {
		t.Fatalf("roundtrip mismatch: %+v", rt)
	}
}

func TestActionCompletedEventJSON(t *testing.T) {
	evt := ActionCompletedEvent{
		GameId: "G1", UnitId: "U1", Phase: "Movement",
		BattlefieldId:  "BF1",
		SourceLocation: map[string]int{"x": 1, "y": 2},
		TargetLocation: map[string]int{"x": 3, "y": 4},
		Unit:           map[string]interface{}{"Id": "U1", "Owner": "CPU"},
	}
	b, err := json.Marshal(evt)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var rt ActionCompletedEvent
	if err := json.Unmarshal(b, &rt); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if rt.GameId != "G1" || rt.UnitId != "U1" || rt.Phase != "Movement" {
		t.Fatalf("roundtrip mismatch: %+v", rt)
	}
}

func TestCloudEventJSON(t *testing.T) {
	ce := CloudEvent{Data: "{}", DataContentType: "application/json", ID: "id1", PubsubName: "redis-pubsub", Source: "gm", SpecVersion: "1.0", Time: "2024-01-01T00:00:00Z", Topic: "t", Type: "x"}
	b, err := json.Marshal(ce)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var rt CloudEvent
	if err := json.Unmarshal(b, &rt); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if rt.ID != ce.ID || rt.Type != ce.Type {
		t.Fatalf("roundtrip mismatch: %+v", rt)
	}
}
