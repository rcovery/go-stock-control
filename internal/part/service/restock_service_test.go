package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/rcovery/go-stock-control/internal/part"
	"github.com/rcovery/go-stock-control/internal/part/repository/memory"
)

func newPart(id, name string, currentStock, minimumStock, averageDailySales int, criticality part.Criticality) part.Part {
	return part.Part{
		ID:                part.ID(id),
		Name:              name,
		CurrentStock:      currentStock,
		MinimumStock:      minimumStock,
		AverageDailySales: averageDailySales,
		CriticalityLevel:  criticality,
	}
}

func TestRestockServiceListPriorities(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewRepository()
	svc := NewRestockService(repo)

	for _, p := range []part.Part{
		newPart("id-alpha", "Alpha", 10, 15, 4, 3),
		newPart("id-beta", "Beta", 10, 13, 4, 5),
		newPart("id-gamma", "Gamma", 10, 13, 6, 5),
		newPart("id-delta", "Delta", 10, 30, 1, 1),
		newPart("id-epsilon", "Epsilon", 10, 13, 6, 5),
		newPart("id-zulu", "Zulu", 100, 20, 4, 5),
	} {
		if err := repo.Create(ctx, p); err != nil {
			t.Fatalf("Create(%s): %v", p.Name, err)
		}
	}

	got, err := svc.ListPriorities(ctx)
	if err != nil {
		t.Fatalf("ListPriorities: %v", err)
	}

	expected := []part.RestockPriority{
		{PartID: "id-delta", Name: "Delta", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 30, UrgencyScore: 20},
		{PartID: "id-epsilon", Name: "Epsilon", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 13, UrgencyScore: 15},
		{PartID: "id-gamma", Name: "Gamma", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 13, UrgencyScore: 15},
		{PartID: "id-beta", Name: "Beta", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 13, UrgencyScore: 15},
		{PartID: "id-alpha", Name: "Alpha", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 15, UrgencyScore: 15},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ListPriorities() =\n%+v\nwant\n%+v", got, expected)
	}
}

type erroringReader struct{}

func (erroringReader) List(ctx context.Context) ([]part.Part, error) {
	return nil, errors.New("list failed")
}

func (erroringReader) ListByCategory(ctx context.Context, category string) ([]part.Part, error) {
	return nil, nil
}

func TestRestockServiceListPrioritiesPropagatesError(t *testing.T) {
	ctx := context.Background()
	svc := NewRestockService(erroringReader{})

	_, err := svc.ListPriorities(ctx)
	if err == nil {
		t.Fatal("ListPriorities with failing reader: want error, got nil")
	}
	if err.Error() != "list failed" {
		t.Errorf("ListPriorities: got error %v, want 'list failed'", err)
	}
}
