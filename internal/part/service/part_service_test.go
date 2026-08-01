package service

import (
	"context"
	"errors"
	"testing"

	"github.com/rcovery/go-stock-control/internal/part"
	"github.com/rcovery/go-stock-control/internal/part/errs"
	"github.com/rcovery/go-stock-control/internal/part/repository/memory"
)

type recordingRepository struct {
	*memory.Repository
	createCalls int
	updateCalls int
	deleteCalls int
}

func (r *recordingRepository) Create(ctx context.Context, p part.Part) error {
	r.createCalls++
	return r.Repository.Create(ctx, p)
}

func (r *recordingRepository) Update(ctx context.Context, p part.Part) error {
	r.updateCalls++
	return r.Repository.Update(ctx, p)
}

func (r *recordingRepository) Delete(ctx context.Context, id part.ID) error {
	r.deleteCalls++
	return r.Repository.Delete(ctx, id)
}

func validPart() part.Part {
	return part.Part{
		Name:              "Filtro de Óleo",
		Category:          "engine",
		CurrentStock:      15,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          18.5,
		CriticalityLevel:  3,
	}
}

func TestPartServiceCreate(t *testing.T) {
	ctx := context.Background()
	svc := NewPartService(memory.NewRepository())

	p := validPart()
	created, err := svc.Create(ctx, p)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if created.ID == "" {
		t.Fatal("Create: ID not assigned")
	}
	if created.Name != p.Name ||
		created.Category != p.Category ||
		created.CurrentStock != p.CurrentStock ||
		created.MinimumStock != p.MinimumStock ||
		created.AverageDailySales != p.AverageDailySales ||
		created.LeadTimeDays != p.LeadTimeDays ||
		created.UnitCost != p.UnitCost ||
		created.CriticalityLevel != p.CriticalityLevel {
		t.Errorf("Create: got %+v, want same data as %+v plus ID", created, p)
	}
}

func TestPartServiceCreateValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &recordingRepository{Repository: memory.NewRepository()}
	svc := NewPartService(repo)

	p := validPart()
	p.CriticalityLevel = part.Criticality(0)

	_, err := svc.Create(ctx, p)
	if err == nil {
		t.Fatal("Create with invalid part: want error, got nil")
	}
	var validationErr errs.ValidationError
	if !errors.As(err, &validationErr) {
		t.Errorf("Create with invalid part: want ValidationError, got %v", err)
	}
	if repo.createCalls != 0 {
		t.Errorf("Create with invalid part: repository called %d times, want 0", repo.createCalls)
	}
}

func TestPartServiceUpdate(t *testing.T) {
	ctx := context.Background()
	svc := NewPartService(memory.NewRepository())

	created, err := svc.Create(ctx, validPart())
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	created.CurrentStock = 42
	updated, err := svc.Update(ctx, created)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.CurrentStock != 42 {
		t.Errorf("Update: CurrentStock = %d, want 42", updated.CurrentStock)
	}

	listed, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(listed) != 1 || listed[0].CurrentStock != 42 {
		t.Errorf("List after update = %+v, want CurrentStock 42", listed)
	}
}

func TestPartServiceUpdateNotFound(t *testing.T) {
	ctx := context.Background()
	svc := NewPartService(memory.NewRepository())

	p := validPart()
	p.ID = part.ID("missing-id")

	_, err := svc.Update(ctx, p)
	if !errors.Is(err, errs.NotFoundError) {
		t.Errorf("Update missing: want NotFoundError, got %v", err)
	}
}

func TestPartServiceUpdateValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &recordingRepository{Repository: memory.NewRepository()}
	svc := NewPartService(repo)

	created, err := svc.Create(ctx, validPart())
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	created.CriticalityLevel = part.Criticality(6)
	_, err = svc.Update(ctx, created)
	if err == nil {
		t.Fatal("Update with invalid part: want error, got nil")
	}
	var validationErr errs.ValidationError
	if !errors.As(err, &validationErr) {
		t.Errorf("Update with invalid part: want ValidationError, got %v", err)
	}
	if repo.updateCalls != 0 {
		t.Errorf("Update with invalid part: repository called %d times, want 0", repo.updateCalls)
	}
}

func TestPartServiceDelete(t *testing.T) {
	ctx := context.Background()
	svc := NewPartService(memory.NewRepository())

	created, err := svc.Create(ctx, validPart())
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	listed, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(listed) != 0 {
		t.Errorf("List after delete = %+v, want empty", listed)
	}
}

func TestPartServiceDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	svc := NewPartService(memory.NewRepository())

	err := svc.Delete(ctx, part.ID("missing-id"))
	if !errors.Is(err, errs.NotFoundError) {
		t.Errorf("Delete missing: want NotFoundError, got %v", err)
	}
}

func TestPartServiceList(t *testing.T) {
	ctx := context.Background()
	svc := NewPartService(memory.NewRepository())

	alpha := validPart()
	alpha.Name = "Alpha"
	alpha.Category = "engine"

	beta := validPart()
	beta.Name = "Beta"
	beta.Category = "transmission"

	if _, err := svc.Create(ctx, alpha); err != nil {
		t.Fatalf("Create(Alpha): %v", err)
	}
	if _, err := svc.Create(ctx, beta); err != nil {
		t.Fatalf("Create(Beta): %v", err)
	}

	all, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("List len = %d, want 2", len(all))
	}

	engines, err := svc.ListByCategory(ctx, "engine")
	if err != nil {
		t.Fatalf("ListByCategory: %v", err)
	}
	if len(engines) != 1 || engines[0].Name != "Alpha" {
		t.Errorf("ListByCategory = %+v, want only Alpha", engines)
	}
}
