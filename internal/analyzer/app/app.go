package app

import (
	"context"
	"log"
	"text_analyzer/internal/analyzer/web/server"

	"go.uber.org/fx"
)

func NewApp(lc fx.Lifecycle, server *server.Server, logger *log.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Start; err != nil {
					logger.Fatalf("Server failed to start: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Println("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}
