package di

import (
	"log"
	"net/http"
	"os"
	"time"

	"text_analyzer/internal/config"
	"text_analyzer/internal/receiver/app"
	"text_analyzer/internal/receiver/client"
	"text_analyzer/internal/receiver/model"
	"text_analyzer/internal/receiver/repository"
	"text_analyzer/internal/receiver/service"
	"text_analyzer/internal/receiver/web/handler"
	"text_analyzer/internal/receiver/web/server"

	"go.uber.org/fx"
)

func newLogger() *log.Logger {
	return log.New(os.Stdout, "[receiver] ", log.LstdFlags|log.Lshortfile)
}

var Module = fx.Module("service-A",
	fx.Provide(
		config.NewConfig,
		newLogger,
		func() *http.Client {
			return &http.Client{Timeout: 5 * time.Second}
		},
		func(cfg *config.Config, clientHTTP *http.Client, logger *log.Logger) service.AnalyzerClient {
			return client.NewAnalyzerClient(cfg.AnalyzerURL, clientHTTP, logger)
		},
		repository.NewTextRepository,
		func(cfg *config.Config, repo model.TextRepository, analyzer service.AnalyzerClient, logger *log.Logger) service.ReceiverService {
			return service.NewReceiverService(repo, analyzer, logger, cfg.AnalysisTimeout)
		},
		handler.NewReceiverHandler,
		server.NewServer,
	),
	fx.Invoke(app.NewApp),
)
