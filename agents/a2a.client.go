package agents

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

/* NOTE:
This A2A protocol implementation is a subset of the A2A specification.
IMPORTANT:
This is a work in progress and may not cover all aspects of the A2A protocol.
*/

func (agent *Agent) PingAgent(agentBaseURL string) (AgentCard, error) {
	resp, err := http.Get(agentBaseURL + "/.well-known/agent.json")
	if err != nil {
		return AgentCard{}, err
	}
	defer resp.Body.Close()

	var agentCard AgentCard
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&agentCard); err != nil {
			return AgentCard{}, err
		}
		return agentCard, nil
	} else {
		return agentCard, errors.New("failed to ping agent: " + resp.Status)
	}
}

func (agent *Agent) SendToAgent(agentBaseURL string, taskRequest TaskRequest) (TaskResponse, error) {
	jsonTaskRequest, err := TaskRequestToJSONString(taskRequest)
	if err != nil {
		return TaskResponse{}, err
	}

	resp, err := http.Post(agentBaseURL+"/", "application/json", strings.NewReader(jsonTaskRequest))
	if err != nil {
		return TaskResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TaskResponse{}, errors.New("failed to send task request: " + resp.Status)
	}

	var taskResponse TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResponse); err != nil {
		return TaskResponse{}, err
	}

	return taskResponse, nil
}
