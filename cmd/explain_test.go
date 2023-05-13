package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Projectorio/dento/cmd/mocks"
	"github.com/franciscoescher/goopenai"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestExplainCmd(t *testing.T) {
	// arrange
	var mockClient = new(mocks.OpenAiClient)
	initTestServer(mockClient, createExpectedRequest())

	buf, cmd := prepareExplainCommand(mockClient)

	// act
	cmd.Execute()

	expectedOutput := getExpectedExplainOutput()
	// assert
	assert.Equal(t, expectedOutput, buf.String())
}

func prepareExplainCommand(mockClient *mocks.OpenAiClient) (*bytes.Buffer, *cobra.Command) {
	buf := bytes.NewBufferString("")
	cmd := NewExplainCommand(mockClient.CreateCompletions)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"Deployment"})
	return buf, cmd
}

func getExpectedExplainOutput() string {
	var expectedOutput = `Deployment

A Deployment provides declarative updates for Pods and ReplicaSets.

This information may not be accurate.

Pod
An overview of Pods.
https://kubernetes.io/docs/concepts/workloads/pods/

ReplicaSet
An overview of ReplicaSets.
https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/

`
	return expectedOutput
}

func initTestServer(mockClient *mocks.OpenAiClient, expectedRequest goopenai.CreateCompletionsRequest) {
	// Mock response
	mockResponse := mockOpenAiResponse(mockClient, expectedRequest)

	// Create a new server to mock the API response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer testServer.Close()
	setEnvironmentVariables(testServer)

}

func setEnvironmentVariables(testServer *httptest.Server) {
	os.Setenv("OPEN_AI_API_KEY", "test-api-key")
	os.Setenv("OPEN_AI_API_ORG", "test-org")
	os.Setenv("OPEN_AI_API_URL", testServer.URL)
}

func createExpectedRequest() goopenai.CreateCompletionsRequest {
	var responseFormat, _ = os.ReadFile("./request/response_format.json")
	var givenSearchTerm = "Deployment"
	expectedRequest := goopenai.CreateCompletionsRequest{
		Model: "gpt-3.5-turbo",
		Messages: []goopenai.Message{
			{
				Role: "user",
				Content: fmt.Sprintf(
					"explain %s in context of kubernetes. in your response make a new line every 80"+
						" charecters. also structure your response in a json with the following format "+string(responseFormat), givenSearchTerm),
			},
		},
		Temperature: 0.2,
	}
	return expectedRequest
}

func mockOpenAiResponse(mockClient *mocks.OpenAiClient, expectedRequest goopenai.CreateCompletionsRequest) []byte {
	mockResponse, err := os.ReadFile("../test/explain_response.json")
	if err != nil {
		panic(err)
	}

	response := goopenai.CreateCompletionsResponse{
		Model: "gpt-3.5-turbo",
		Choices: []struct {
			Message struct {
				Role    string `json:"role,omitempty"`
				Content string `json:"content,omitempty"`
			} `json:"message"`
			Text         string      `json:"text,omitempty"`
			Index        int         `json:"index,omitempty"`
			Logprobs     interface{} `json:"logprobs,omitempty"`
			FinishReason string      `json:"finish_reason,omitempty"`
		}{
			{
				Message: goopenai.Message{
					Content: string(mockResponse),
				},
			},
		},
	}
	mockClient.On("CreateCompletions", context.Background(), expectedRequest).Return(
		response,
		nil,
	)
	return mockResponse
}
