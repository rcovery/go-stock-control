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
