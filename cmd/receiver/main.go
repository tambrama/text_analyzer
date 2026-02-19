package main

import (
	"text_analyzer/internal/receiver/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		di.Module,
	).Run()
}
