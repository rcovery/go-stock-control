package part

import (
	"fmt"

	"github.com/rcovery/go-stock-control/internal/part/errs"
)

func (p Part) Validate() error {
	if !p.CriticalityLevel.IsValid() {
		return errs.NewValidationError(fmt.Sprintf(
			"criticalityLevel must be between %d and %d",
			MinCriticality, MaxCriticality,
		))
	}

	return nil
}
