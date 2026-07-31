package part

import "context"

type Reader interface {
	List(ctx context.Context) ([]Part, error)
	ListByCategory(ctx context.Context, category string) ([]Part, error)
}

type Writer interface {
	Create(ctx context.Context, part Part) error
	Update(ctx context.Context, part Part) error
	Delete(ctx context.Context, id ID) error
}

type Repository interface {
	Reader
	Writer
}
