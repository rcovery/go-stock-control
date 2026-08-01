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
	needing := make([]restockCandidate, 0, len(parts))
	for _, p := range parts {
		r := CalculateRestock(p)
		if r.NeedsRestock() {
			needing = append(needing, restockCandidate{part: p, restock: r})
		}
	}

	sortRestockCandidates(needing)

	priorities := make([]RestockPriority, 0, len(needing))
	for _, c := range needing {
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
