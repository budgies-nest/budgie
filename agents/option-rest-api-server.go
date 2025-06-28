package agents

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/openai/openai-go"
)

type HTTPServerConfig struct {
	Port           string
	Endpoint       string
	StreamEndPoint string
}

func GetBytesBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}
// WithHTTPServer configures the HTTP server for the agent.
// It sets up the server with default endpoints and port if not provided.
// The server handles POST requests for chat completions and streams.
// It also provides a mechanism to cancel ongoing completions via a DELETE request (commented out).
// The server uses the provided HTTPServerConfig to configure its behavior.
// If the StreamEndPoint or Endpoint is not specified, it defaults to "/api/chat-stream" and "/api/chat" respectively.
// The default port is set to "8888" if not specified.
func WithHTTPServer(httpServerConfig HTTPServerConfig) AgentOption {
	return func(agent *Agent) {
		agent.httpServerConfig = httpServerConfig

		if httpServerConfig.StreamEndPoint == "" {
			// Default endpoint path for HTTP server
			agent.httpServerConfig.StreamEndPoint = "/api/chat-stream"
		}

		if httpServerConfig.Endpoint == "" {
			// Default endpoint path for HTTP server
			agent.httpServerConfig.Endpoint = "/api/chat"
		}

		if httpServerConfig.Port == "" {
			// Default port for HTTP server
			agent.httpServerConfig.Port = "8888"
		}

		// Create HTTP server
		agent.httpServer = http.NewServeMux()

		shouldIStopTheCompletion := false

		agent.httpServer.HandleFunc("POST "+agent.httpServerConfig.StreamEndPoint, func(response http.ResponseWriter, request *http.Request) {

			// add a flusher
			flusher, ok := response.(http.Flusher)
			if !ok {
				response.Write([]byte("Error: expected http.ResponseWriter to be an http.Flusher"))
			}
			body := GetBytesBody(request)
			// unmarshal the json data
			var data map[string]string
			err := json.Unmarshal(body, &data)
			if err != nil {
				response.Write([]byte("Error: " + err.Error()))
			}

			//systemContent := data["system"]
			userContent := data["user"]

			agent.Params.Messages = append(
				agent.Params.Messages, openai.UserMessage(userContent),
			)

			answer, err := agent.ChatCompletionStream(context.Background(), func(self *Agent, content string, err error) error {
				response.Write([]byte(content))

				flusher.Flush()
				if !shouldIStopTheCompletion {
					return nil
				} else {
					return errors.New("cancelling request")
				}
				//return nil
			})
			if err != nil {
				response.Write([]byte("Error: " + err.Error()))
				return
			}

			agent.Params.Messages = append(
				agent.Params.Messages, openai.AssistantMessage(answer),
			)

		})

		/* TODO:
		// Cancel/Stop the generation of the completion
		mux.HandleFunc("DELETE /api/completion/cancel", func(response http.ResponseWriter, request *http.Request) {
			shouldIStopTheCompletion = true
			response.Write([]byte("cancelling request..."))
		})
		*/

		agent.httpServer.HandleFunc("POST "+agent.httpServerConfig.Endpoint, func(response http.ResponseWriter, request *http.Request) {


			body := GetBytesBody(request)
			// unmarshal the json data
			var data map[string]string
			err := json.Unmarshal(body, &data)
			if err != nil {
				response.Write([]byte("Error: " + err.Error()))
			}

			//systemContent := data["system"]
			userContent := data["user"]

			agent.Params.Messages = append(
				agent.Params.Messages, openai.UserMessage(userContent),
			)

			answer, err := agent.ChatCompletion(context.Background())
			if err != nil {
				response.Write([]byte("Error: " + err.Error()))
				return
			}
			response.Write([]byte(answer))

			agent.Params.Messages = append(
				agent.Params.Messages, openai.AssistantMessage(answer),
			)
		})

	}
}

// QUESTION: how to handle the conversational memory
// TODO: https, authentication token, etc.
