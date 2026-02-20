package main

import (
	"text_analyzer/internal/analyzer/di"

	"go.uber.org/fx"
)

// @title Text Analyzer API (Analyzer Service)
// @version 1.0
// @description Микросервис для анализа текста: подсчет слов, символов, предложений и средней длины слова.
// @description Этот сервис используется внутренне сервисом Receiver, но также предоставляет прямой API для отладки.
// @host localhost:8082
// @BasePath /api/v1
func main() {
	fx.New(
		di.Module,
	).Run()
}
