package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"text_analyzer/internal/common"
	"text_analyzer/internal/receiver/service"
	"time"
)

type analyzerClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *log.Logger
}

func NewAnalyzerClient(baseURL string, httpClient *http.Client, logger *log.Logger) service.AnalyzerClient {
	base := strings.TrimRight(baseURL, "/")
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}
	return &analyzerClient{
		baseURL:    base,
		httpClient: httpClient,
		logger:     logger,
	}
}


func (c *analyzerClient) Analyze(ctx context.Context, text string) (common.TextStats, error) {
	body, err := json.Marshal(common.AnalyzeRequest{Text: text})
	if err != nil {
		return common.TextStats{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/analyze", bytes.NewReader(body))
	if err != nil {
		return common.TextStats{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return common.TextStats{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return common.TextStats{}, fmt.Errorf("analyzer returned status %d", resp.StatusCode)
	}

	var stats common.TextStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return common.TextStats{}, err
	}
	return stats, nil
}
