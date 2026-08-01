package service

import (
	"context"

	"github.com/rcovery/go-stock-control/internal/part"
)

type PartService struct {
	repo part.Repository
}

func NewPartService(repo part.Repository) *PartService {
	return &PartService{
		repo: repo,
	}
}

func (s *PartService) Create(ctx context.Context, p part.Part) (part.Part, error) {
	if err := p.ValidateForCreation(); err != nil {
		return part.Part{}, err
	}

	newID, err := part.NewID()
	if err != nil {
		return part.Part{}, err
	}
	p.ID = newID

	err = s.repo.Create(ctx, p)
	if err != nil {
		return part.Part{}, err
	}

	return p, nil
}

func (s *PartService) List(ctx context.Context) ([]part.Part, error) {
	return s.repo.List(ctx)
}

func (s *PartService) ListByCategory(ctx context.Context, category string) ([]part.Part, error) {
	return s.repo.ListByCategory(ctx, category)
}

func (s *PartService) Update(ctx context.Context, p part.Part) (part.Part, error) {
	if err := p.Validate(); err != nil {
		return part.Part{}, err
	}

	err := s.repo.Update(ctx, p)
	if err != nil {
		return part.Part{}, err
	}

	return p, nil
}

func (s *PartService) Delete(ctx context.Context, id part.ID) error {
	return s.repo.Delete(ctx, id)
}
