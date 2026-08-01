package part

import "testing"

func TestCalculateRestock(t *testing.T) {
	tests := []struct {
		name          string
		part          Part
		expectedStock int
		expectedScore int
		needsRestock  bool
	}{
		{
			name: "needs restock",
			part: Part{
				CurrentStock:      15,
				MinimumStock:      20,
				AverageDailySales: 4,
				LeadTimeDays:      5,
				CriticalityLevel:  3,
			},
			expectedStock: -5,
			expectedScore: 75,
			needsRestock:  true,
		},
		{
			name: "no restock needed",
			part: Part{
				CurrentStock:      50,
				MinimumStock:      20,
				AverageDailySales: 4,
				LeadTimeDays:      5,
				CriticalityLevel:  3,
			},
			expectedStock: 30,
			expectedScore: 0,
			needsRestock:  false,
		},
		{
			name: "exactly at minimum stock",
			part: Part{
				CurrentStock:      40,
				MinimumStock:      20,
				AverageDailySales: 4,
				LeadTimeDays:      5,
				CriticalityLevel:  2,
			},
			expectedStock: 20,
			expectedScore: 0,
			needsRestock:  false,
		},
		{
			name: "zero daily sales",
			part: Part{
				CurrentStock:      10,
				MinimumStock:      20,
				AverageDailySales: 0,
				LeadTimeDays:      5,
				CriticalityLevel:  1,
			},
			expectedStock: 10,
			expectedScore: 10,
			needsRestock:  true,
		},
		{
			name: "zero lead time",
			part: Part{
				CurrentStock:      10,
				MinimumStock:      20,
				AverageDailySales: 4,
				LeadTimeDays:      0,
				CriticalityLevel:  2,
			},
			expectedStock: 10,
			expectedScore: 20,
			needsRestock:  true,
		},
		{
			name: "higher criticality multiplies score",
			part: Part{
				CurrentStock:      10,
				MinimumStock:      20,
				AverageDailySales: 4,
				LeadTimeDays:      5,
				CriticalityLevel:  5,
			},
			expectedStock: -10,
			expectedScore: 150,
			needsRestock:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateRestock(tt.part)

			if got.ProjectedStock != tt.expectedStock {
				t.Errorf("ProjectedStock = %d, want %d", got.ProjectedStock, tt.expectedStock)
			}
			if got.UrgencyScore != tt.expectedScore {
				t.Errorf("UrgencyScore = %d, want %d", got.UrgencyScore, tt.expectedScore)
			}
			if got.NeedsRestock() != tt.needsRestock {
				t.Errorf("NeedsRestock() = %v, want %v", got.NeedsRestock(), tt.needsRestock)
			}
		})
	}
}

func TestSortByRestockPriority(t *testing.T) {
	parts := []Part{
		{Name: "Alpha", UrgencyScore: 10, CriticalityLevel: 3, AverageDailySales: 4},
		{Name: "Beta", UrgencyScore: 10, CriticalityLevel: 5, AverageDailySales: 4},
		{Name: "Gamma", UrgencyScore: 10, CriticalityLevel: 5, AverageDailySales: 6},
		{Name: "Delta", UrgencyScore: 20, CriticalityLevel: 1, AverageDailySales: 1},
		{Name: "Epsilon", UrgencyScore: 10, CriticalityLevel: 5, AverageDailySales: 6},
	}

	SortByRestockPriority(parts)

	expected := []string{"Delta", "Epsilon", "Gamma", "Beta", "Alpha"}
	for i, want := range expected {
		if parts[i].Name != want {
			t.Errorf("position %d = %s, want %s", i, parts[i].Name, want)
		}
	}
}
