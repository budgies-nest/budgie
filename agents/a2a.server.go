package agents

/* NOTE:
	This A2A protocol implementation is a subset of the A2A specification.
	IMPORTANT:
	This is a work in progress and may not cover all aspects of the A2A protocol.
*/

import "net/http"

func (agent *Agent) StartA2AServer() error {
	errListening := http.ListenAndServe(":"+agent.a2aServerConfig.Port, agent.a2aServer)
	if errListening != nil {
		return errListening
	}
	return nil
}

// A2AServer returns the A2A server mux
func (agent *Agent) A2AServer() *http.ServeMux {
	return agent.a2aServer
}


func (agent *Agent) A2AServerConfig() A2AServerConfig {
	return agent.a2aServerConfig
}