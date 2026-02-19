package server

import (
	"context"
	"log"
	"net/http"
	"text_analyzer/internal/config"
	"text_analyzer/internal/middleware"
	"text_analyzer/internal/receiver/web/handler"
	"time"
)

type Server struct {
	cfg     *config.Config
	server  *http.Server
	handler *handler.ReceiverHandler
	logger  *log.Logger
}

func NewServer(cfg *config.Config, handler *handler.ReceiverHandler, logger *log.Logger) *Server {
	return &Server{cfg: cfg, handler: handler, logger: logger}
}

func (s *Server) Start() error {
	mux := s.handler.Router()
	middleware := middleware.Logging(mux)

	s.server = &http.Server{
		Addr:         ":" + s.cfg.ServerPort,
		Handler:      middleware,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	s.logger.Printf("Сервер запущен на порту %s", s.cfg.ServerPort)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log.Println("Завершение работы сервера...")
	return s.server.Shutdown(shutdownCtx)
}
