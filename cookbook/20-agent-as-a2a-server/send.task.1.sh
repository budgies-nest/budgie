#!/bin/bash
: <<'COMMENT'
# Send task to the agent server
COMMENT

HTTP_PORT=8888
AGENT_BASE_URL=http://0.0.0.0:${HTTP_PORT}

# host.docker.internal

read -r -d '' DATA <<- EOM
{
    "jsonrpc": "2.0",
    "id": "1111",
    "method": "message/send",
    "params": {
      "message": {
        "role": "user",
        "parts": [
          {
            "text": "What is the best pizza in the world?"
          }
        ]
      },
      "metadata": {
        "skill": "ask_for_something"
      }
    }
}
EOM

curl ${AGENT_BASE_URL} \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d "${DATA}" | jq '.'


