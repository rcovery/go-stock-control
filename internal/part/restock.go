package part

import "sort"

type Restock struct {
	ProjectedStock int
	UrgencyScore   int
}

func CalculateRestock(p Part) Restock {
	expectedConsumption := p.AverageDailySales * p.LeadTimeDays
	projectedStock := p.CurrentStock - expectedConsumption

	var urgencyScore int
	if projectedStock < p.MinimumStock {
		urgencyScore = (p.MinimumStock - projectedStock) * p.CriticalityLevel
	}

	return Restock{
		ProjectedStock: projectedStock,
		UrgencyScore:   urgencyScore,
	}
}

func (r Restock) NeedsRestock() bool {
	return r.UrgencyScore > 0
}

func SortByRestockPriority(parts []Part) {
	sort.SliceStable(parts, func(i, j int) bool {
		a, b := parts[i], parts[j]
		if a.UrgencyScore != b.UrgencyScore {
			return a.UrgencyScore > b.UrgencyScore
		}
		if a.CriticalityLevel != b.CriticalityLevel {
			return a.CriticalityLevel > b.CriticalityLevel
		}
		if a.AverageDailySales != b.AverageDailySales {
			return a.AverageDailySales > b.AverageDailySales
		}
		return a.Name < b.Name
	})
}

type RestockPriority struct {
	PartID         ID     `json:"partId"`
	Name           string `json:"name"`
	CurrentStock   int    `json:"currentStock"`
	ProjectedStock int    `json:"projectedStock"`
	MinimumStock   int    `json:"minimumStock"`
	UrgencyScore   int    `json:"urgencyScore"`
}

func BuildRestockPriorities(parts []Part) []RestockPriority {
	needing := make([]Part, 0, len(parts))
	for _, p := range parts {
		r := CalculateRestock(p)
		p.ProjectedStock = r.ProjectedStock
		p.UrgencyScore = r.UrgencyScore

		if r.NeedsRestock() {
			needing = append(needing, p)
		}
	}

	SortByRestockPriority(needing)

	priorities := make([]RestockPriority, 0, len(needing))
	for _, p := range needing {
		priorities = append(priorities, RestockPriority{
			PartID:         p.ID,
			Name:           p.Name,
			CurrentStock:   p.CurrentStock,
			ProjectedStock: p.ProjectedStock,
			MinimumStock:   p.MinimumStock,
			UrgencyScore:   p.UrgencyScore,
		})
	}

	return priorities
}
