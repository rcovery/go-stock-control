package part

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
