package agents

import (
	"net/http"
)

func (agent *Agent) StartHttpServer() error {

	errListening := http.ListenAndServe(":"+agent.httpServerConfig.Port, agent.httpServer)

	return errListening
}

func (agent *Agent) HttpServer() *http.ServeMux {
	return agent.httpServer
}
// TODO: perhaps add a helper more straightforward
