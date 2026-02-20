package service

import (
	"log"
	"strings"
	"unicode"

	"text_analyzer/internal/common"
)

type AnalyzerService interface {
	Analyze(req string) common.TextStats
}

type analyzerService struct {
	logger *log.Logger
}

func NewAnalyzerService(logger *log.Logger) AnalyzerService {
	return &analyzerService{logger: logger}
}

func (s *analyzerService) Analyze(req string) common.TextStats {
	text := strings.TrimSpace(req)
	words := splitWords(text)
	sentences := splitSentences(text)

	var totalLetters int
	for _, w := range words {
		for _, r := range w {
			if unicode.IsLetter(r) {
				totalLetters++
			}
		}
	}

	avgLen := 0.0
	if len(words) > 0 {
		avgLen = float64(totalLetters) / float64(len(words))
	}

	return common.TextStats{
		WordsCount:      len(words),
		CharactersCount: len([]rune(text)),
		SentencesCount:  len(sentences),
		AverageWordLen:  avgLen,
	}
}

// splitWords возвращает слайт слов
func splitWords(text string) []string {
	//возвращаем слайс строк между разделителями
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || (unicode.IsPunct(r) && r != '-' && r != '\'')
	})

	return fields
}

// splitSentences возвращает слайс предложений
func splitSentences(text string) []string {
	separators := ".!?"
	var sentences []string
	start := 0
	runes := []rune(text)
	for i, r := range runes {
		if strings.ContainsRune(separators, r) {
			s := strings.TrimSpace(string(runes[start : i+1]))
			if s != "" {
				sentences = append(sentences, s)
			}
			start = i + 1
		}
	}
	if start < len(runes) {
		s := strings.TrimSpace(string(runes[start:]))
		if s != "" {
			sentences = append(sentences, s)
		}
	}
	return sentences
}
