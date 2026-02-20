package handler

import (
	"encoding/json"
	"net/http"
	_ "text_analyzer/docs"
	"text_analyzer/internal/analyzer/service"
	"text_analyzer/internal/common"

	httpSwagger "github.com/swaggo/http-swagger"
)

type AnalyzerHandler struct {
	svc service.AnalyzerService
}

func NewAnalyzerHandler(svc service.AnalyzerService) *AnalyzerHandler {
	return &AnalyzerHandler{svc: svc}
}

func (h *AnalyzerHandler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/analyze", h.handleAnalyze)
	mux.HandleFunc("/api/v1/health", h.handleHealth)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return mux
}

// handleAnalyze godoc
// @Summary      Анализировать текст
// @Description  Принимает текст, подсчитывает статистику (слова, символы, предложения) и возвращает результат немедленно.
// @Tags         analyzer
// @Accept       json
// @Produce      json
// @Param        request  body      common.AnalyzeRequest  true  "Объект запроса с полем 'text'"
// @Success      200      {object}  common.TextStats       "Статистика текста"
// @Failure      400      {string}  string                 "Некорректный JSON или пустой текст"
// @Failure      500      {string}  string                 "Внутренняя ошибка сервера при анализе"
// @Router       /api/v1/analyze [post]
func (h *AnalyzerHandler) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req common.AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	stats := h.svc.Analyze(req.Text)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stats)
}

// handleHealth godoc
// @Summary      Проверка здоровья сервиса анализа
// @Description  Эндпоинт для мониторинга доступности сервиса анализа текста.
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string "Сервис доступен"
// @Router       /api/v1/health [get]
func (h *AnalyzerHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
