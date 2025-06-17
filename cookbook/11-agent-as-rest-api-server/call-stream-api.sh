#!/bin/bash

curl --no-buffer http://localhost:8080/api/chat-stream \
-H "Content-Type: application/json" \
-d '
{
  "user": "who is Jean-Luc Picard?"
}' 
