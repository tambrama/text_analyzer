package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	analyzerService "text_analyzer/internal/analyzer/service"
	analyzerHandler "text_analyzer/internal/analyzer/web/handler"
	"text_analyzer/internal/common"
	"text_analyzer/internal/receiver/client"
	"text_analyzer/internal/receiver/model"
	"text_analyzer/internal/receiver/repository"
	"text_analyzer/internal/receiver/service"
	"text_analyzer/internal/receiver/web/handler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullFlow(t *testing.T) {
	//заглушка
	logger := log.New(&testWriter{t}, "[INTEGRATION] ", log.LstdFlags)

	//сервис Б
	aService := analyzerService.NewAnalyzerService(logger)
	aHandler := analyzerHandler.NewAnalyzerHandler(aService)
	aRouter := aHandler.Router()
	serverB := httptest.NewServer(aRouter)
	defer serverB.Close()

	t.Logf("Analyzer service started at: %s", serverB.URL)

	//сервис А
	repo := repository.NewTextRepository()

	//клиент для отправки данных в сервис Б
	httpClient := &http.Client{Timeout: 5 * time.Second}
	aClient := client.NewAnalyzerClient(serverB.URL, httpClient, logger)

	rService := service.NewReceiverService(repo, aClient, logger, 2*time.Second)
	rHandler := handler.NewReceiverHandler(rService, logger)
	rRouter := rHandler.Router()

	serverA := httptest.NewServer(rRouter)
	defer serverA.Close()

	t.Logf("Receiver service started at: %s", serverA.URL)

	// отправляем текст в Сервис А
	textRequest := common.AnalyzeRequest{Text: "Hello world! This is a test."}
	body, err := json.Marshal(textRequest)
	require.NoError(t, err)

	resp, err := http.Post(serverA.URL+"/api/v1/text", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusAccepted, resp.StatusCode)

	var submitResp struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&submitResp)
	require.NoError(t, err)
	require.NotEmpty(t, submitResp.ID)
	t.Logf("Task created with ID: %s", submitResp.ID)

	// ждем обработки, тк обработка асинхронная
	time.Sleep(1 * time.Second)

	statusURL := fmt.Sprintf("%s/api/v1/status/%s", serverA.URL, submitResp.ID)
	statusResp, err := http.Get(statusURL)
	require.NoError(t, err)
	defer statusResp.Body.Close()

	assert.Equal(t, http.StatusOK, statusResp.StatusCode)

	var taskResp model.Text
	err = json.NewDecoder(statusResp.Body).Decode(&taskResp)
	require.NoError(t, err)

	assert.Equal(t, "DONE", string(taskResp.Status))
	assert.Empty(t, taskResp.Error, "Task should not have an error")

	assert.Equal(t, 6, taskResp.Result.WordsCount)
	assert.Greater(t, taskResp.Result.CharactersCount, 20)
	assert.Greater(t, taskResp.Result.SentencesCount, 0)

	t.Logf("Analysis result: %+v", taskResp.Result)
}

type testWriter struct {
	t *testing.T
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(string(p))
	return len(p), nil
}
