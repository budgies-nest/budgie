package agents

/* NOTE:
	This A2A protocol implementation is a subset of the A2A specification.
	IMPORTANT:
	This is a work in progress and may not cover all aspects of the A2A protocol.
*/

import "encoding/json"

func TaskRequestToJSONString(taskRequest TaskRequest) (string, error) {
	jsonData, err := json.MarshalIndent(taskRequest, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func AgentCardToJSONString(agentCard AgentCard) (string, error) {
	jsonData, err := json.MarshalIndent(agentCard, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func TaskResponseToJSONString(taskResponse TaskResponse) (string, error) {
	jsonData, err := json.MarshalIndent(taskResponse, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
