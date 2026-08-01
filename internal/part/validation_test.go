package part

import "testing"

func TestPartValidateCriticalityLevel(t *testing.T) {
	for _, level := range []Criticality{MinCriticality, 3, MaxCriticality} {
		p := Part{CriticalityLevel: level}
		if err := p.Validate(); err != nil {
			t.Errorf("CriticalityLevel %d should be valid, got error: %v", level, err)
		}
	}

	for _, level := range []Criticality{MinCriticality - 1, MaxCriticality + 1, 0, -5, 100} {
		p := Part{CriticalityLevel: level}
		if err := p.Validate(); err == nil {
			t.Errorf("CriticalityLevel %d should be invalid", level)
		}
	}
}
