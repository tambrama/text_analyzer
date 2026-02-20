package service_test

import (
	"testing"
	"text_analyzer/internal/analyzer/service"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeService_CountStats(t *testing.T) {
	svc := service.NewAnalyzerService(nil)

	tests := []struct {
		name      string
		input     string
		wantWords int
	}{
		{"simple", "Hello world", 2},
		{"with punctuation", "Hello, world!", 2},
		{"empty", "", 0},
		{"single", "Go", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := svc.Analyze(tt.input)
			assert.Equal(t, tt.wantWords, stats.WordsCount)
		})
	}
}
