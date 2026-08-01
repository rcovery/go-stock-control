package part

import (
	"fmt"
	"strings"

	"github.com/rcovery/go-stock-control/internal/part/errs"
)

func (p Part) Validate() error {
	if !p.CriticalityLevel.IsValid() {
		return errs.NewValidationError(fmt.Sprintf(
			"criticalityLevel must be between %d and %d",
			MinCriticality, MaxCriticality,
		))
	}

	if strings.TrimSpace(p.Name) == "" {
		return errs.NewValidationError("name must not be empty")
	}

	if p.AverageDailySales < 0 {
		return errs.NewValidationError("averageDailySales must not be negative")
	}

	if p.LeadTimeDays < 0 {
		return errs.NewValidationError("leadTimeDays must not be negative")
	}

	if p.UnitCost < 0 {
		return errs.NewValidationError("unitCost must not be negative")
	}

	if p.MinimumStock < 0 {
		return errs.NewValidationError("minimumStock must not be negative")
	}

	return nil
}
