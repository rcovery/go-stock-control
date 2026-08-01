package part

import "testing"

func TestPartValidateCriticalityLevel(t *testing.T) {
	for _, level := range []Criticality{MinCriticality, 3, MaxCriticality} {
		p := Part{Name: "Filtro", CriticalityLevel: level}
		if err := p.Validate(); err != nil {
			t.Errorf("CriticalityLevel %d should be valid, got error: %v", level, err)
		}
	}

	for _, level := range []Criticality{MinCriticality - 1, MaxCriticality + 1, 0, -5, 100} {
		p := Part{Name: "Filtro", CriticalityLevel: level}
		if err := p.Validate(); err == nil {
			t.Errorf("CriticalityLevel %d should be invalid", level)
		}
	}
}

func TestPartValidateName(t *testing.T) {
	for _, name := range []string{"", " ", "\t\n", "   "} {
		p := Part{Name: name, CriticalityLevel: 3}
		if err := p.Validate(); err == nil {
			t.Errorf("Name %q should be invalid", name)
		}
	}

	p := Part{Name: "  Filtro de Óleo  ", CriticalityLevel: 3}
	if err := p.Validate(); err != nil {
		t.Errorf("Name with surrounding spaces should be valid, got error: %v", err)
	}
}

func TestPartValidateNonNegativeFields(t *testing.T) {
	valid := Part{
		Name:              "Filtro",
		CurrentStock:      10,
		MinimumStock:      5,
		AverageDailySales: 2,
		LeadTimeDays:      3,
		UnitCost:          1.5,
		CriticalityLevel:  3,
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("fully valid part should pass Validate(), got: %v", err)
	}

	invalidCases := []struct {
		name   string
		mutate func(*Part)
	}{
		{"AverageDailySales", func(p *Part) { p.AverageDailySales = -1 }},
		{"LeadTimeDays", func(p *Part) { p.LeadTimeDays = -1 }},
		{"UnitCost", func(p *Part) { p.UnitCost = -1 }},
		{"MinimumStock", func(p *Part) { p.MinimumStock = -1 }},
	}
	for _, tc := invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			p := valid
			tc.mutate(&p)
			if err := p.Validate(); err == nil {
				t.Errorf("%s negative should be invalid", tc.name)
			}
		})
	}
}

func TestPartValidateForCreationCurrentStock(t *testing.T) {
	p := Part{
		Name:              "Filtro",
		CurrentStock:      -5,
		MinimumStock:      10,
		AverageDailySales: 2,
		LeadTimeDays:      3,
		UnitCost:          1.5,
		CriticalityLevel:  3,
	}

	if err := p.Validate(); err != nil {
		t.Errorf("Validate() should allow negative CurrentStock, got: %v", err)
	}

	if err := p.ValidateForCreation(); err == nil {
		t.Error("ValidateForCreation() should reject negative CurrentStock, got nil")
	}
}
