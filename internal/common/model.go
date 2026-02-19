package common

// AnalyzeRequest запрос на анализ текста
// @Description Содержит текстовую строку для подсчета статистики (слова, символы, предложения)
// swagger:model AnalyzeRequest
type AnalyzeRequest struct {
	Text string `json:"text" example:"Hello world! This is a sample text for analysis." validate:"required,min=1"`
}

// TextStats результат анализа текста
// @Description Статистические данные о тексте: количество слов, символов, предложений и средняя длина слова
// swagger:model TextStats
type TextStats struct {
	WordsCount      int     `json:"words_count" example:"10"`
	CharactersCount int     `json:"characters_count" example:"55"`
	SentencesCount  int     `json:"sentences_count" example:"2"`
	AverageWordLen  float64 `json:"average_word_length" example:"4.5"`
}