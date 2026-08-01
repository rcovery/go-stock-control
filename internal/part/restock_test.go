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
		{
			name: "negative current stock",
			part: Part{
				CurrentStock:      -5,
				MinimumStock:      20,
				AverageDailySales: 4,
				LeadTimeDays:      5,
				CriticalityLevel:  2,
			},
			expectedStock: -25,
			expectedScore: 90,
			needsRestock:  true,
		},
		{
			name: "high lead time",
			part: Part{
				CurrentStock:      100,
				MinimumStock:      20,
				AverageDailySales: 10,
				LeadTimeDays:      30,
				CriticalityLevel:  3,
			},
			expectedStock: -200,
			expectedScore: 660,
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
