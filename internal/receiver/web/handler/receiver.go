package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"text_analyzer/internal/receiver/model"
	"text_analyzer/internal/receiver/service"
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
	return mux
}

// handleSubmitText godoc
// @Summary      Отправить текст на анализ
// @Description  Принимает текстовую строку, инициирует процесс анализа в фоновом режиме и возвращает уникальный ID задачи.
// @Tags         receiver
// @Accept       json
// @Produce      json
// @Param        request  body      TextRequest  true  "Объект запроса, содержащий поле 'text'"
// @Success      202      {object}  SubmitResponse  "Задача принята в обработку. Поле 'id' содержит идентификатор для проверки статуса."
// @Failure      400      {object}  ErrorResponse   "Некорректный JSON или пустой текст"
// @Failure      500      {object}  ErrorResponse   "Внутренняя ошибка сервера при создании задачи"
// @Router       /api/v1/text [post]
func (h *ReceiverHandler) handleSubmitText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("submit decode error: %v", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	resp, err := h.service.SubmitText(r.Context(), req)
	if err != nil {
		if err == model.ErrValidation {
			http.Error(w, "text is required", http.StatusBadRequest)
			return
		}
		h.logger.Printf("submit error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

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
// @Success      200  {object}  model.TextRequest "Объект задачи со статусом и результатами"
// @Failure      400  {object}  ErrorResponse   "Некорректный или пустой ID"
// @Failure      404  {object}  ErrorResponse   "Задача с указанным ID не найдена"
// @Failure      500  {object}  ErrorResponse   "Внутренняя ошибка сервера"
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

	tr, err := h.service.GetStatus(r.Context(), id)
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
	_ = json.NewEncoder(w).Encode(tr)
}

// handleHealth godoc
// @Summary      Проверка здоровья сервиса (Health Check)
// @Description  Эндпоинт для мониторинга доступности сервиса. Возвращает статус "ok", если сервис работает корректно.
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse "Сервис доступен"
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
