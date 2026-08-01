package part

import (
	"reflect"
	"testing"
)

func TestBuildRestockPriorities(t *testing.T) {
	parts := []Part{
		{ID: "id-alpha", Name: "Alpha", CurrentStock: 10, MinimumStock: 15, AverageDailySales: 4, CriticalityLevel: 3},
		{ID: "id-beta", Name: "Beta", CurrentStock: 10, MinimumStock: 13, AverageDailySales: 4, CriticalityLevel: 5},
		{ID: "id-gamma", Name: "Gamma", CurrentStock: 10, MinimumStock: 13, AverageDailySales: 6, CriticalityLevel: 5},
		{ID: "id-delta", Name: "Delta", CurrentStock: 10, MinimumStock: 30, AverageDailySales: 1, CriticalityLevel: 1},
		{ID: "id-epsilon", Name: "Epsilon", CurrentStock: 10, MinimumStock: 13, AverageDailySales: 6, CriticalityLevel: 5},
		{ID: "id-zulu", Name: "Zulu", CurrentStock: 100, MinimumStock: 20, AverageDailySales: 4, CriticalityLevel: 5},
	}

	got := BuildRestockPriorities(parts)

	expected := []RestockPriority{
		{PartID: "id-delta", Name: "Delta", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 30, UrgencyScore: 20},
		{PartID: "id-epsilon", Name: "Epsilon", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 13, UrgencyScore: 15},
		{PartID: "id-gamma", Name: "Gamma", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 13, UrgencyScore: 15},
		{PartID: "id-beta", Name: "Beta", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 13, UrgencyScore: 15},
		{PartID: "id-alpha", Name: "Alpha", CurrentStock: 10, ProjectedStock: 10, MinimumStock: 15, UrgencyScore: 15},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("BuildRestockPriorities() =\n%+v\nwant\n%+v", got, expected)
	}
}
