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
		urgencyScore = (p.MinimumStock - projectedStock) * int(p.CriticalityLevel)
	}

	return Restock{
		ProjectedStock: projectedStock,
		UrgencyScore:   urgencyScore,
	}
}

func (r Restock) NeedsRestock() bool {
	return r.UrgencyScore > 0
}

type restockCandidate struct {
	part    Part
	restock Restock
}

func sortRestockCandidates(candidates []restockCandidate) {
	sort.SliceStable(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.restock.UrgencyScore != b.restock.UrgencyScore {
			return a.restock.UrgencyScore > b.restock.UrgencyScore
		}
		if a.part.CriticalityLevel != b.part.CriticalityLevel {
			return a.part.CriticalityLevel > b.part.CriticalityLevel
		}
		if a.part.AverageDailySales != b.part.AverageDailySales {
			return a.part.AverageDailySales > b.part.AverageDailySales
		}
		return a.part.Name < b.part.Name
	})
}

func SortByRestockPriority(parts []Part) {
	candidates := make([]restockCandidate, 0, len(parts))
	for _, p := range parts {
		candidates = append(candidates, restockCandidate{part: p, restock: CalculateRestock(p)})
	}

	sortRestockCandidates(candidates)

	for i, c := range candidates {
		parts[i] = c.part
	}
}
