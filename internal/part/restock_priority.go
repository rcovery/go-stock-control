package part

import "sort"

type RestockPriority struct {
	PartID         ID     `json:"partId"`
	Name           string `json:"name"`
	CurrentStock   int    `json:"currentStock"`
	ProjectedStock int    `json:"projectedStock"`
	MinimumStock   int    `json:"minimumStock"`
	UrgencyScore   int    `json:"urgencyScore"`
}

type restockCandidate struct {
	part    Part
	restock Restock
}

func BuildRestockPriorities(parts []Part) []RestockPriority {
	candidates := make([]restockCandidate, 0, len(parts))
	for _, p := range parts {
		r := CalculateRestock(p)
		if r.NeedsRestock() {
			candidates = append(candidates, restockCandidate{part: p, restock: r})
		}
	}

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

	priorities := make([]RestockPriority, 0, len(candidates))
	for _, c := range candidates {
		priorities = append(priorities, RestockPriority{
			PartID:         c.part.ID,
			Name:           c.part.Name,
			CurrentStock:   c.part.CurrentStock,
			ProjectedStock: c.restock.ProjectedStock,
			MinimumStock:   c.part.MinimumStock,
			UrgencyScore:   c.restock.UrgencyScore,
		})
	}

	return priorities
}
