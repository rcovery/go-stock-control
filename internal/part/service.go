package part

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, p Part) (Part, error) {
	newID, err := NewID()
	if err != nil {
		return Part{}, err
	}
	p.ID = newID

	r := CalculateRestock(p)
	p.ProjectedStock = r.ProjectedStock
	p.UrgencyScore = r.UrgencyScore

	err = s.repo.Create(ctx, p)
	if err != nil {
		return Part{}, err
	}

	return p, nil
}

func (s *Service) List(ctx context.Context) ([]Part, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListByCategory(ctx context.Context, category string) ([]Part, error) {
	return s.repo.ListByCategory(ctx, category)
}

func (s *Service) Update(ctx context.Context, p Part) error {
	r := CalculateRestock(p)
	p.ProjectedStock = r.ProjectedStock
	p.UrgencyScore = r.UrgencyScore

	return s.repo.Update(ctx, p)
}

func (s *Service) Delete(ctx context.Context, id ID) error {
	return s.repo.Delete(ctx, id)
}
