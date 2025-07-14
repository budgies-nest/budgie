package agents

/* NOTE:
	This A2A protocol implementation is a subset of the A2A specification.
	IMPORTANT:
	This is a work in progress and may not cover all aspects of the A2A protocol.
*/

// AgentCard represents the metadata for this agent
type AgentCard struct {
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	URL          string           `json:"url"`
	Version      string           `json:"version"`
	Capabilities map[string]any   `json:"capabilities"`
	Skills       []map[string]any `json:"skills,omitempty"` // Optional, for storing skills related to the agent
}

type AgentMessageParams struct {
	Message  AgentMessage   `json:"message"`
	MetaData map[string]any `json:"metadata,omitempty"` // Optional, for additional metadata
}

// REF: https://google-a2a.github.io/A2A/specification/#92-basic-execution-synchronous-polling-style
// TaskRequest represents an incoming A2A task request
type TaskRequest struct {
	JSONRpcVersion string             `json:"jsonrpc"` // Should be "2.0"
	ID             string             `json:"id"`
	Params         AgentMessageParams `json:"params"`
	Method         string             `json:"method,omitempty"` // Optional, for specifying the method of the task
}

// Message represents a message structure
type AgentMessage struct {
	Role      string     `json:"role,omitempty"`
	Parts     []TextPart `json:"parts"`
	MessageID string     `json:"messageId,omitempty"` // Optional, for storing message ID
	TaskID    string     `json:"taskId,omitempty"`    // Optional, for storing task ID
	ContextID string     `json:"contextId,omitempty"` // Optional, for storing context ID
}

// TextPart represents a text part of a message
type TextPart struct {
	Text string `json:"text"`
	Type string `json:"type"` // Should be "text" for text parts
}

// TaskStatus represents the status of a task
type TaskStatus struct {
	State string `json:"state"`
}

// TODO: make the response compliant with the A2A protocol
// REF: https://google-a2a.github.io/A2A/specification/#92-basic-execution-synchronous-polling-style

type Artifact struct {
	ArtifactID string     `json:"artifactId"`
	Name       string     `json:"name"`
	Parts      []TextPart `json:"parts"` // Parts of the artifact, e.g., text, images, etc.
}

type Result struct {
	ID        string         `json:"id"`
	ContextID string         `json:"contextId"`
	Status    TaskStatus     `json:"status"`
	Artifacts []Artifact     `json:"artifacts,omitempty"` // Optional, for storing artifacts related to the task
	History   []AgentMessage `json:"history,omitempty"`   // Optional, for storing message history related to the task
	Kind      string         `json:"kind"`                // Should be "task"
	Metadata  map[string]any `json:"metadata,omitempty"`  // Optional, for additional metadata
}

// TaskResponse represents the response task structure
type TaskResponse struct {
	JSONRpcVersion string `json:"jsonrpc"` // Should be "2.0"
	ID             string `json:"id"`
	Result         Result `json:"result"` // The result of the task execution

}
