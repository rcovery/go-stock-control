package part

// Define o que é uma peça.
// Seria possível transformar o nome e a categoria em value objects também
type Part struct {
	ID                ID      `json:"id"`
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	CurrentStock      int     `json:"currentStock"`
	MinimumStock      int     `json:"minimumStock"`
	AverageDailySales int     `json:"averageDailySales"`
	LeadTimeDays      int     `json:"leadTimeDays"`
	UnitCost          float64 `json:"unitCost"`
	CriticalityLevel  int     `json:"criticalityLevel"`
	ProjectedStock    int     `json:"-"`
	UrgencyScore      int     `json:"-"`
}
