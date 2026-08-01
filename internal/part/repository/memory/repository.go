package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/rcovery/go-stock-control/internal/part"
	"github.com/rcovery/go-stock-control/internal/part/errs"
)

// Repository stores parts in memory. The map is not safe for concurrent
// access, and HTTP handlers run in one goroutine per request, so the
// RWMutex guards all reads (RLock) and writes (Lock).
type Repository struct {
	mu    sync.RWMutex
	parts map[part.ID]part.Part
}

func NewRepository() *Repository {
	return &Repository{
		parts: make(map[part.ID]part.Part),
	}
}

func (r *Repository) List(ctx context.Context) ([]part.Part, error) {
	r.mu.RLock() // concurrent reads allowed
	defer r.mu.RUnlock()

	parts := make([]part.Part, 0, len(r.parts))
	for _, p := range r.parts {
		parts = append(parts, p)
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})

	return parts, nil
}

func (r *Repository) ListByCategory(ctx context.Context, category string) ([]part.Part, error) {
	r.mu.RLock() // concurrent reads allowed
	defer r.mu.RUnlock()

	var parts []part.Part
	for _, p := range r.parts {
		if p.Category == category {
			parts = append(parts, p)
		}
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})

	return parts, nil
}

func (r *Repository) Create(ctx context.Context, p part.Part) error {
	r.mu.Lock() // exclusive write access
	defer r.mu.Unlock()

	if _, ok := r.parts[p.ID]; ok {
		return errs.NotCreatedErr.New(fmt.Sprintf("Create: part %q already exists", p.ID))
	}

	r.parts[p.ID] = p

	return nil
}

func (r *Repository) Update(ctx context.Context, p part.Part) error {
	r.mu.Lock() // exclusive write access
	defer r.mu.Unlock()

	if _, ok := r.parts[p.ID]; !ok {
		return errs.NotFoundError.New(fmt.Sprintf("Update: part %q not found", p.ID))
	}

	r.parts[p.ID] = p

	return nil
}

func (r *Repository) Delete(ctx context.Context, id part.ID) error {
	r.mu.Lock() // exclusive write access
	defer r.mu.Unlock()

	if _, ok := r.parts[id]; !ok {
		return errs.NotFoundError.New(fmt.Sprintf("Delete: part %q not found", id))
	}

	delete(r.parts, id)

	return nil
}
