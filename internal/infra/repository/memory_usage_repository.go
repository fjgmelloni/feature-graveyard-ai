package repository

import (
	"context"
	"sync"

	"feature-graveyard-ai/internal/domain/feature"
)

type MemoryUsageRepository struct {
	mu   sync.RWMutex
	logs []feature.UsageLog
}

func NewMemoryUsageRepository(seed []feature.UsageLog) *MemoryUsageRepository {
	return &MemoryUsageRepository{logs: append([]feature.UsageLog(nil), seed...)}
}

func (r *MemoryUsageRepository) SaveMany(_ context.Context, logs []feature.UsageLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logs = append(r.logs, logs...)
	return nil
}

func (r *MemoryUsageRepository) List(_ context.Context) ([]feature.UsageLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return append([]feature.UsageLog(nil), r.logs...), nil
}
