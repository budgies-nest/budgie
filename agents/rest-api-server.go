package agents

import (
	"net/http"
)

func (agent *Agent) StartHttpServer() error {

	errListening := http.ListenAndServe(":"+agent.httpServerConfig.Port, agent.httpServer)

	return errListening
}
