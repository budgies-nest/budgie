package agents

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/openai/openai-go"
)

// go test -v -run TestLogger
func TestLogger(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	
	// Create a logger that writes to our buffer
	logger := &Logger{
		level:   LogLevelDebug,
		logger:  log.New(&buf, "", 0),
		enabled: true,
	}

	// Test logging a chat completion
	request := openai.ChatCompletionNewParams{
		Model: "test-model",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("test message"),
		},
	}

	logger.LogChatCompletion("test-agent", request, "test response", 100*time.Millisecond, nil)

	// Check that log was written
	logOutput := buf.String()
	if logOutput == "" {
		t.Error("Expected log output, but got empty string")
	}

	// Parse the JSON log entry
	var entry LogEntry
	err := json.Unmarshal([]byte(logOutput), &entry)
	if err != nil {
		t.Errorf("Failed to parse log entry as JSON: %v", err)
	}

	// Verify log entry contents
	if entry.Type != "chat_completion" {
		t.Errorf("Expected type 'chat_completion', got '%s'", entry.Type)
	}
	if entry.AgentName != "test-agent" {
		t.Errorf("Expected agent name 'test-agent', got '%s'", entry.AgentName)
	}
	if entry.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", entry.Level)
	}
}

func TestLoggerEnabled(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:   LogLevelInfo,
		logger:  log.New(&buf, "", 0),
		enabled: false, // Disabled
	}

	request := openai.ChatCompletionNewParams{
		Model: "test-model",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("test message"),
		},
	}

	logger.LogChatCompletion("test-agent", request, "test response", 100*time.Millisecond, nil)

	// Should be no output when disabled
	if buf.String() != "" {
		t.Error("Expected no log output when logger is disabled, but got output")
	}
}

func TestLoggerLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:   LogLevelError, // Only error level
		logger:  log.New(&buf, "", 0),
		enabled: true,
	}

	request := openai.ChatCompletionNewParams{
		Model: "test-model",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("test message"),
		},
	}

	// This should not log (info level when logger is set to error level)
	logger.LogChatCompletion("test-agent", request, "test response", 100*time.Millisecond, nil)

	if buf.String() != "" {
		t.Error("Expected no log output for info level when logger is set to error level")
	}
}

func TestGlobalLogger(t *testing.T) {
	// Test global logger functions
	originalLogger := GetGlobalLogger()
	
	// Test enabling logging
	EnableLogging(LogLevelInfo)
	if !GetGlobalLogger().IsEnabled() {
		t.Error("Expected global logger to be enabled after EnableLogging")
	}

	// Test disabling logging
	DisableLogging()
	if GetGlobalLogger().IsEnabled() {
		t.Error("Expected global logger to be disabled after DisableLogging")
	}

	// Restore original logger
	SetGlobalLogger(originalLogger)
}