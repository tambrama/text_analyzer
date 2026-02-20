package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "text_analyzer/docs"
	"text_analyzer/internal/common"
	"text_analyzer/internal/receiver/model"
	"text_analyzer/internal/receiver/service"

	httpSwagger "github.com/swaggo/http-swagger"
)

type ReceiverHandler struct {
	service service.ReceiverService
	logger  *log.Logger
}

func NewReceiverHandler(service service.ReceiverService, logger *log.Logger) *ReceiverHandler {
	return &ReceiverHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ReceiverHandler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/text", h.handleSubmitText)
	mux.HandleFunc("/api/v1/status/", h.handleGetStatus)
	mux.HandleFunc("/api/v1/health", h.handleHealth)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return mux
}

// handleSubmitText godoc
// @Summary      Отправить текст на анализ
// @Description  Принимает текстовую строку, инициирует процесс анализа в фоновом режиме и возвращает уникальный ID задачи.
// @Tags         receiver
// @Accept       json
// @Produce      json
// @Param        request  body      common.AnalyzeRequest  true  "Объект запроса, содержащий поле 'text'"
// @Failure      400  {string}  string               "Некорректный или пустой ID"
// @Failure      404  {string}  string               "Задача с указанным ID не найдена"
// @Failure      500  {string}  string               "Внутренняя ошибка сервера"
// @Router       /api/v1/text [post]
func (h *ReceiverHandler) handleSubmitText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req common.AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("submit decode error: %v", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := h.service.SubmitText(r.Context(), req.Text)
	if err != nil {
		if err == model.ErrValidation {
			http.Error(w, "text is required", http.StatusBadRequest)
			return
		}
		h.logger.Printf("submit error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	resp := model.SubmitResponse{ID: id}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(resp)
}

// handleGetStatus godoc
// @Summary      Получить статус задачи анализа
// @Description  Возвращает текущий статус обработки текста и результаты анализа, если задача завершена.
// @Tags         receiver
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Уникальный идентификатор задачи (UUID)"
// @Success      200  {object}  model.Text "Объект задачи со статусом и результатами"
// @Failure      400  {string}  string               "Некорректный или пустой ID"
// @Failure      404  {string}  string               "Задача с указанным ID не найдена"
// @Failure      500  {string}  string               "Внутренняя ошибка сервера"
// @Router       /api/v1/status/{id} [get]
func (h *ReceiverHandler) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/status/")
	id := strings.TrimSpace(path)
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetStatus(r.Context(), id)
	if err != nil {
		if err == model.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		h.logger.Printf("get status error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(status)
}

// handleHealth godoc
// @Summary      Проверка здоровья сервиса (Health Check)
// @Description  Эндпоинт для мониторинга доступности сервиса. Возвращает статус "ok", если сервис работает корректно.
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  string "Сервис доступен"
// @Router       /api/v1/health [get]
func (h *ReceiverHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
