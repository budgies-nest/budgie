package agents

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/openai/openai-go"
)

type LogLevel int

const (
	LogLevelOff LogLevel = iota
	LogLevelError
	LogLevelInfo
	LogLevelDebug
)

type Logger struct {
	level   LogLevel
	logger  *log.Logger
	enabled bool
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Type      string                 `json:"type"`
	AgentName string                 `json:"agent_name,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

func NewLogger(level LogLevel, enabled bool) *Logger {
	return &Logger{
		level:   level,
		logger:  log.New(os.Stdout, "", 0),
		enabled: enabled,
	}
}

func (l *Logger) SetEnabled(enabled bool) {
	l.enabled = enabled
}

func (l *Logger) IsEnabled() bool {
	return l.enabled
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) logEntry(entry LogEntry) {
	if !l.enabled || l.level == LogLevelOff {
		return
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return
	}
	l.logger.Println(string(jsonData))
}

func (l *Logger) LogChatCompletion(agentName string, request openai.ChatCompletionNewParams, response string, duration time.Duration, err error) {
	if !l.enabled || l.level < LogLevelInfo {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "info",
		Type:      "chat_completion",
		AgentName: agentName,
		Data: map[string]interface{}{
			"model":             request.Model,
			"messages_count":    len(request.Messages),
			"max_tokens":        request.MaxTokens,
			"temperature":       request.Temperature,
			"response_length":   len(response),
			"duration_ms":       duration.Milliseconds(),
		},
	}

	if err != nil {
		entry.Level = "error"
		entry.Error = err.Error()
	} else {
		entry.Message = "Chat completion successful"
	}

	l.logEntry(entry)
}

func (l *Logger) LogChatCompletionStream(agentName string, request openai.ChatCompletionNewParams, response string, duration time.Duration, err error) {
	if !l.enabled || l.level < LogLevelInfo {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "info",
		Type:      "chat_completion_stream",
		AgentName: agentName,
		Data: map[string]interface{}{
			"model":             request.Model,
			"messages_count":    len(request.Messages),
			"max_tokens":        request.MaxTokens,
			"temperature":       request.Temperature,
			"response_length":   len(response),
			"duration_ms":       duration.Milliseconds(),
		},
	}

	if err != nil {
		entry.Level = "error"
		entry.Error = err.Error()
	} else {
		entry.Message = "Chat completion stream successful"
	}

	l.logEntry(entry)
}

func (l *Logger) LogToolsCompletion(agentName string, request openai.ChatCompletionNewParams, toolCalls []openai.ChatCompletionMessageToolCall, duration time.Duration, err error) {
	if !l.enabled || l.level < LogLevelInfo {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "info",
		Type:      "tools_completion",
		AgentName: agentName,
		Data: map[string]interface{}{
			"model":            request.Model,
			"messages_count":   len(request.Messages),
			"tools_count":      len(request.Tools),
			"tool_calls_count": len(toolCalls),
			"duration_ms":      duration.Milliseconds(),
		},
	}

	if err != nil {
		entry.Level = "error"
		entry.Error = err.Error()
	} else {
		entry.Message = fmt.Sprintf("Tools completion successful with %d tool calls", len(toolCalls))
		if l.level >= LogLevelDebug {
			toolNames := make([]string, len(toolCalls))
			for i, tc := range toolCalls {
				toolNames[i] = tc.Function.Name
			}
			entry.Data["tool_names"] = toolNames
		}
	}

	l.logEntry(entry)
}

func (l *Logger) LogToolExecution(agentName string, toolName string, args map[string]any, response string, duration time.Duration, err error) {
	if !l.enabled || l.level < LogLevelInfo {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "info",
		Type:      "tool_execution",
		AgentName: agentName,
		Data: map[string]interface{}{
			"tool_name":       toolName,
			"response_length": len(response),
			"duration_ms":     duration.Milliseconds(),
		},
	}

	if l.level >= LogLevelDebug {
		entry.Data["args"] = args
		entry.Data["response"] = response
	}

	if err != nil {
		entry.Level = "error"
		entry.Error = err.Error()
	} else {
		entry.Message = fmt.Sprintf("Tool '%s' executed successfully", toolName)
	}

	l.logEntry(entry)
}

func (l *Logger) LogMCPToolExecution(agentName string, toolName string, args map[string]any, response string, clientType string, duration time.Duration, err error) {
	if !l.enabled || l.level < LogLevelInfo {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "info",
		Type:      "mcp_tool_execution",
		AgentName: agentName,
		Data: map[string]interface{}{
			"tool_name":       toolName,
			"client_type":     clientType,
			"response_length": len(response),
			"duration_ms":     duration.Milliseconds(),
		},
	}

	if l.level >= LogLevelDebug {
		entry.Data["args"] = args
		entry.Data["response"] = response
	}

	if err != nil {
		entry.Level = "error"
		entry.Error = err.Error()
	} else {
		entry.Message = fmt.Sprintf("MCP tool '%s' executed successfully via %s", toolName, clientType)
	}

	l.logEntry(entry)
}

func (l *Logger) LogError(agentName string, errorType string, message string, err error, context map[string]interface{}) {
	if !l.enabled || l.level < LogLevelError {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "error",
		Type:      errorType,
		AgentName: agentName,
		Message:   message,
		Error:     err.Error(),
		Data:      context,
	}

	l.logEntry(entry)
}

var defaultLogger = NewLogger(LogLevelOff, false)

func SetGlobalLogger(logger *Logger) {
	defaultLogger = logger
}

func GetGlobalLogger() *Logger {
	return defaultLogger
}

func EnableLogging(level LogLevel) {
	defaultLogger.SetLevel(level)
	defaultLogger.SetEnabled(true)
}

func DisableLogging() {
	defaultLogger.SetEnabled(false)
}