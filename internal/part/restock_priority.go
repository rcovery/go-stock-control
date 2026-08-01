package part

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
