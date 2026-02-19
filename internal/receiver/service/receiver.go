package service

import (
	"context"
	"errors"
	"io"
	"log"
	"text_analyzer/internal/common"
	"text_analyzer/internal/receiver/model"
	"time"

	"github.com/google/uuid"
)

// AnalyzerClient посредник между работой сервисов
type AnalyzerClient interface {
	Analyze(ctx context.Context, text string) (common.TextStats, error)
}

type ReceiverService interface {
	SubmitText(ctx context.Context, text string) (string, error)
	GetStatus(ctx context.Context, id string) (model.Text, error)
}

type receiverService struct {
	repo            model.TextRepository
	analyzer        AnalyzerClient
	logger          *log.Logger
	analysisTimeout time.Duration
}

func NewReceiverService(repo model.TextRepository, analyzer AnalyzerClient, logger *log.Logger, analysisTimeout time.Duration) ReceiverService {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	if analysisTimeout == 0 {
		analysisTimeout = 3 * time.Second
	}
	return &receiverService{
		repo:            repo,
		analyzer:        analyzer,
		logger:          logger,
		analysisTimeout: analysisTimeout,
	}
}

// SubmitText валидирует текст, отправляет в анализатор
func (u *receiverService) SubmitText(ctx context.Context, text string) (string, error) {
	if len(text) == 0 {
		return "", model.ErrValidation
	}
	if u.analyzer == nil {
		return "", errors.New("analyzer client is not configured")
	}

	id := uuid.New().String()

	task := model.Text{
		ID:     id,
		Text:   text,
		Status: model.StatusPending,
	}

	if err := u.repo.Save(ctx, task); err != nil {
		return "", err
	}

	go u.runAnalysis(context.Background(), id, task)

	return id, nil
}

// runAnalysis запускает анализ текста с таймаутом
func (u *receiverService) runAnalysis(ctx context.Context, id string, text model.Text) {
	//таймаут на анализ 
	ctx, cancel := context.WithTimeout(ctx, u.analysisTimeout)
	defer cancel()
	
	stats, err := u.analyzer.Analyze(ctx, text.Text)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Printf("analyzer timeout: %v", err)
			u.markFailed(id, err.Error(), model.StatusTimeout)
			return
		}
		u.logger.Printf("analyzer error: %v", err)
		u.markFailed(id, err.Error(), model.StatusFailed)
		return
	}

	if err := u.updateWithResult(id, stats); err != nil {
		u.logger.Printf("update result error: %v", err)
	}
}
// markFailed обновляет статус на "провал"
func (u *receiverService) markFailed(id, msg string, status model.WorkStatus) {
	ctx := context.Background()
	tr, err := u.repo.Get(ctx, id)
	if err != nil {
		u.logger.Printf("failed to get task", "id", id, "err", err)
		return
	}
	tr.Status = status
	tr.Error = msg
	 if err := u.repo.Update(ctx, tr); err != nil {
        u.logger.Printf("failed to update task", "id", id, "err", err)
		return
    }
}
// updateWithResult обновляет результат
func (u *receiverService) updateWithResult(id string, stats common.TextStats) error {
	ctx := context.Background()
	tr, err := u.repo.Get(ctx, id)
	if err != nil {
		u.logger.Printf("failed to get task", "id", id, "err", err)
		return err
	}
	tr.Status = model.StatusDone
	tr.Result = &stats
	return u.repo.Update(ctx, tr)
}

// GetStatus возвращает статус обработки
func (u *receiverService) GetStatus(ctx context.Context, id string) (model.Text, error) {
	return u.repo.Get(ctx, id)
}
