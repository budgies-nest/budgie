#!/bin/bash

curl --no-buffer http://localhost:5050/api/fight \
-H "Content-Type: application/json" \
-d '
{
  "info": "✋ INFO: 12 12 68"
}' 
