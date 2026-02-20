package main

import (
	"text_analyzer/internal/receiver/di"

	"go.uber.org/fx"
)

// @title Text Analyzer API (Receiver Service)
// @version 1.0
// @description REST API для приема текстов на анализ и проверки статуса задач.
// @description Сервис принимает текст, возвращает ID задачи и позволяет проверить статус обработки.
// @host localhost:8081
// @BasePath /api/v1

func main() {
	fx.New(
		di.Module,
	).Run()
}
