#!/bin/bash

curl --no-buffer http://localhost:5050/api/info \
-H "Content-Type: application/json" \
-d '
{
  "info": "✋ INFO: 0123456789"
}' 
