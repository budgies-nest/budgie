package agents

/* NOTE:
	This A2A protocol implementation is a subset of the A2A specification.
	IMPORTANT:
	This is a work in progress and may not cover all aspects of the A2A protocol.
*/

import (
	"encoding/json"
	"log"
	"net/http"
)

type A2AServerConfig struct {
	Port string
	//Endpoint       string
	//StreamEndPoint string
}

func WithA2AServer(a2aServerConfig A2AServerConfig) AgentOption {
	return func(agent *Agent) {
		agent.a2aServerConfig = a2aServerConfig
		agent.a2aServer = http.NewServeMux()

		// Register handlers
		agent.a2aServer.HandleFunc("/.well-known/agent.json", agent.getAgentCard)

		// Use the synchronous handler instead
		agent.a2aServer.HandleFunc("/", agent.handleTaskSync)

		//agent.a2aServer.HandleFunc(a2aServerConfig.Endpoint, agent.handleA2ATaskRequest)
		//agent.a2aServer.HandleFunc(a2aServerConfig.StreamEndPoint, agent.handleA2ATaskStream)
	}
}

// Serve the Agent Card at the well-known URL
func (agent *Agent) getAgentCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agent.agentCard)
}

// Alternative synchronous implementation that should work better
func (agent *Agent) handleTaskSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var taskRequest TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
		http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
		return
	}

	switch taskRequest.Method {
	case "message/send":
		if len(taskRequest.Params.Message.Parts) > 0 {
			// Process the task synchronously without mutex in the HTTP handler
			// The mutex should only be in the AgentCallback if needed
			responseTask, err := agent.agentCallback(taskRequest)
			if err != nil {
				log.Printf("Agent callback failed for task %s: %v", taskRequest.ID, err)
				http.Error(w, `{"error": "agent callback failed"}`, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(responseTask)
		} else {
			http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, `{"error": "unknown method"}`, http.StatusBadRequest)
	}
}


func WithAgentCard(agentCard AgentCard) AgentOption {
	return func(agent *Agent) {
		agent.agentCard = agentCard
	}
}

func WithAgentCallback(callback func(ctx *AgentCallbackContext) (TaskResponse, error)) AgentOption {
	return func(agent *Agent) {
		agent.agentCallback = func(taskRequest TaskRequest) (TaskResponse, error) {
			ctx := &AgentCallbackContext{
				CompletionContext: CompletionContext{
					Agent:   agent,
					Context: nil, // You can pass a proper context.Context if available
				},
				TaskRequest:  &taskRequest,
				TaskResponse: nil,
			}
			return callback(ctx)
		}
	}
}
