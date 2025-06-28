#!/bin/bash

curl http://localhost:8080/api/chat \
-H "Content-Type: application/json" \
-d '
{
  "user": "who is James T Kirk?"
}' 