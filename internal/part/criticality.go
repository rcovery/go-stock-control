package part

type Criticality int

const (
	MinCriticality Criticality = 1
	MaxCriticality Criticality = 5
)

func (c Criticality) IsValid() bool {
	return c >= MinCriticality && c <= MaxCriticality
}
