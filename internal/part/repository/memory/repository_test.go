package memory

import (
	"context"
	"testing"

	"github.com/rcovery/go-stock-control/internal/part"
	"github.com/rcovery/go-stock-control/internal/part/errs"
)

func newPart(id, name, category string) part.Part {
	return part.Part{
		ID:       part.ID(id),
		Name:     name,
		Category: category,
	}
}

func TestRepositoryCRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository()

	zeta := newPart("1", "Zeta", "transmissao")
	alpha := newPart("2", "Alpha", "motor")

	if err := repo.Create(ctx, zeta); err != nil {
		t.Fatalf("Create(zeta): %v", err)
	}
	if err := repo.Create(ctx, alpha); err != nil {
		t.Fatalf("Create(alpha): %v", err)
	}

	created, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(created) != 2 {
		t.Fatalf("List len = %d, want 2", len(created))
	}
	if created[0].Name != "Alpha" || created[1].Name != "Zeta" {
		t.Fatalf("List not ordered by name: %+v", created)
	}

	motors, err := repo.ListByCategory(ctx, "motor")
	if err != nil {
		t.Fatalf("ListByCategory: %v", err)
	}
	if len(motors) != 1 || motors[0].Name != "Alpha" {
		t.Fatalf("ListByCategory = %+v, want only Alpha", motors)
	}

	zeta.CurrentStock = 42
	if err := repo.Update(ctx, zeta); err != nil {
		t.Fatalf("Update: %v", err)
	}

	updated, err := repo.ListByCategory(ctx, "transmissao")
	if err != nil {
		t.Fatalf("ListByCategory: %v", err)
	}
	if len(updated) != 1 || updated[0].CurrentStock != 42 {
		t.Fatalf("updated = %+v, want CurrentStock 42", updated)
	}

	if err := repo.Delete(ctx, zeta.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	remaining, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(remaining) != 1 {
		t.Fatalf("after delete len = %d, want 1", len(remaining))
	}
}

func TestRepositoryDuplicateCreate(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository()

	p := newPart("1", "Alpha", "motor")
	if err := repo.Create(ctx, p); err != nil {
		t.Fatalf("Create: %v", err)
	}

	err := repo.Create(ctx, p)
	if err == nil {
		t.Fatal("Create duplicate: want error, got nil")
	}
}

func TestRepositoryMissingID(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository()

	p := newPart("1", "Alpha", "motor")

	if err := repo.Update(ctx, p); !errs.NotFoundError.Is(err) {
		t.Fatalf("Update missing: want NotFoundError, got %v", err)
	}

	if err := repo.Delete(ctx, p.ID); !errs.NotFoundError.Is(err) {
		t.Fatalf("Delete missing: want NotFoundError, got %v", err)
	}
}
