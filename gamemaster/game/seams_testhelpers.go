package game

import (
	"context"

	"github.com/corstijank/mekstrike/domain/unit"
	bfclient "github.com/corstijank/mekstrike/gamemaster/clients/bfclient"
	uclient "github.com/corstijank/mekstrike/gamemaster/clients/uclient"
	dclient "github.com/dapr/go-sdk/client"
)

// Helper to stub battlefield client factory during tests.
func SetBattlefieldClientFactoryForTest(f func(dclient.Client, string) *bfclient.BattlefieldClient) (restore func()) {
	prev := getBattlefieldClient
	getBattlefieldClient = f
	return func() { getBattlefieldClient = prev }
}

// Helper to stub unit client factory during tests.
func SetUnitClientFactoryForTest(f func(dclient.Client, string, string, *unit.Stats) *uclient.UnitClient) (restore func()) {
	prev := newUnitClient
	newUnitClient = f
	return func() { newUnitClient = prev }
}

// Fake army builder used in tests.
type fakeAB struct{ out []*unit.Stats }

func (f *fakeAB) CreateArmy(_ context.Context, _, _, _, _ int) ([]*unit.Stats, error) {
	return f.out, nil
}

// Use a trivial armybuilder that always returns the provided units.
func UseArmyBuilderStubForTest(units []*unit.Stats) (restore func()) {
	prev := newArmyBuilder
	newArmyBuilder = func(_ dclient.Client) (armyBuilder, error) { return &fakeAB{out: units}, nil }
	return func() { newArmyBuilder = prev }
}
