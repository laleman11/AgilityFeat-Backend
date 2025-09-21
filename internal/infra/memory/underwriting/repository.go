package underwriting

import (
	"context"
	"sync"

	"AgilityFeat-Backend/internal/port"
)

type Repository struct {
	mu      sync.RWMutex
	records map[string][]port.EvaluationRecord
}

func NewRepository() *Repository {
	return &Repository{records: make(map[string][]port.EvaluationRecord)}
}

func (r *Repository) Save(ctx context.Context, record port.EvaluationRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.records[record.UserID] = append(r.records[record.UserID], record)
	return nil
}

func (r *Repository) FindByUser(ctx context.Context, userID string) ([]port.EvaluationRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := r.records[userID]
	out := make([]port.EvaluationRecord, len(entries))
	copy(out, entries)
	return out, nil
}
