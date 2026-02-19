package di

import (
	"log"
	"os"

	"text_analyzer/internal/analyzer/app"
	"text_analyzer/internal/analyzer/web/handler"
	"text_analyzer/internal/analyzer/web/server"
	"text_analyzer/internal/analyzer/service"
	"text_analyzer/internal/config"

	"go.uber.org/fx"
)

func newLogger() *log.Logger {
	return log.New(os.Stdout, "[receiver] ", log.LstdFlags|log.Lshortfile)
}

var Module = fx.Module("service-A",
	fx.Provide(
		config.NewConfig,
		newLogger,
		service.NewAnalyzerService,
		handler.NewAnalyzerHandler,
		server.NewServer,
	),
	fx.Invoke(app.NewApp),
)
