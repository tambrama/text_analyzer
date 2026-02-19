package repository

import (
	"context"
	"sync"
	"text_analyzer/internal/receiver/model"
	"time"
)

type textRepository struct {
	mu   sync.RWMutex
	data map[string]model.Text
}

func NewTextRepository() model.TextRepository {
	return &textRepository{
		data: make(map[string]model.Text),
	}
}

func (r *textRepository) Save(ctx context.Context, text model.Text) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	text.CreatedAt = now
	text.UpdatedAt = now
	r.data[text.ID] = text
	return nil
}

func (r *textRepository) Get(ctx context.Context, id string) (model.Text, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.data[id]; ok {
		return v, nil
	}
	return model.Text{}, model.ErrNotFound
}

func (r *textRepository) Update(ctx context.Context, text model.Text) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[text.ID]; !ok {
		return model.ErrNotFound
	}
	text.UpdatedAt = time.Now()
	r.data[text.ID] = text
	return nil
}
