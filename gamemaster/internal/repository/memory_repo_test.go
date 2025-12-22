package repository

import (
	"context"
	"testing"
	"time"

	"github.com/corstijank/mekstrike/gamemaster/game"
)

func TestMemoryRepo_CRUDAndExists(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryGameRepository()

	g := &game.Data{ID: "G1", StartTime: time.Now()}
	if err := r.Save(ctx, g); err != nil {
		t.Fatalf("save error: %v", err)
	}

	ok, err := r.Exists(ctx, "G1")
	if err != nil || !ok {
		t.Fatalf("exists = %v, err=%v, want true,nil", ok, err)
	}

	got, err := r.Get(ctx, "G1")
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	if got.ID != "G1" {
		t.Fatalf("got ID = %s, want G1", got.ID)
	}

	if err := r.Delete(ctx, "G1"); err != nil {
		t.Fatalf("delete error: %v", err)
	}
	ok, _ = r.Exists(ctx, "G1")
	if ok {
		t.Fatalf("exists after delete = %v, want false", ok)
	}
}

func TestMemoryRepo_GetAllSortedByStartTime(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryGameRepository()

	t1 := time.Now().Add(-time.Hour)
	t2 := time.Now()
	t3 := time.Now().Add(time.Hour)

	_ = r.Save(ctx, &game.Data{ID: "G3", StartTime: t3})
	_ = r.Save(ctx, &game.Data{ID: "G1", StartTime: t1})
	_ = r.Save(ctx, &game.Data{ID: "G2", StartTime: t2})

	all, err := r.GetAll(ctx)
	if err != nil {
		t.Fatalf("getall error: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("len = %d, want 3", len(all))
	}
	if all[0].ID != "G1" || all[1].ID != "G2" || all[2].ID != "G3" {
		t.Fatalf("order = [%s,%s,%s], want [G1,G2,G3]", all[0].ID, all[1].ID, all[2].ID)
	}
}
