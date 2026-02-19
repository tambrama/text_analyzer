package model

import (
	"context"
	"text_analyzer/internal/common"
	"time"
)

// WorkStatus статус обработки задачи анализа
// @Description Возможные состояния задачи в процессе обработки текста
// swagger:enum WorkStatus
type WorkStatus string

const (
	StatusPending WorkStatus = "PENDING" //в ожидании
	StatusDone    WorkStatus = "DONE"    //готово
	StatusFailed  WorkStatus = "FAILED"  //сбой
	StatusTimeout WorkStatus = "TIMEOUT"
)


// Text запрос и результат анализа текста
// @Description Основная сущность, содержащая исходный текст, текущий статус и результаты статистики
// swagger:model Text
type Text struct {
	ID        string            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Text      string            `json:"text" example:"Hello world! This is a test sentence."`
	Status    WorkStatus        `json:"status" example:"DONE"`
	CreatedAt time.Time         `json:"created_at" example:"2026-02-19T10:00:00Z"`
	UpdatedAt time.Time         `json:"updated_at" example:"2026-02-19T10:00:05Z"`
	Error     string            `json:"error,omitempty" example:"context deadline exceeded"`
	Result    *common.TextStats `json:"result,omitempty"`
}

// TextRepository интерфейс хранилища задач
// @Description Определяет методы для сохранения, получения и обновления задач анализа

type TextRepository interface {
	Save(ctx context.Context, text Text) error
	Get(ctx context.Context, id string) (Text, error)
	Update(ctx context.Context, text Text) error
}
