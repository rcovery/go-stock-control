package service

import (
	"context"

	"github.com/rcovery/go-stock-control/internal/part"
)

type RestockService struct {
	reader part.Reader
}

func NewRestockService(reader part.Reader) *RestockService {
	return &RestockService{
		reader: reader,
	}
}

func (s *RestockService) ListPriorities(ctx context.Context) ([]part.RestockPriority, error) {
	parts, err := s.reader.List(ctx)
	if err != nil {
		return nil, err
	}

	return part.BuildRestockPriorities(parts), nil
}
