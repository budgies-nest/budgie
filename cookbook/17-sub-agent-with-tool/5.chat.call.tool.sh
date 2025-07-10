#!/bin/bash

curl --no-buffer http://localhost:5050/api/chat-stream \
-H "Content-Type: application/json" \
-d '
{
  "user": "Say hello to Bob. Who is Spock?"
}' 
