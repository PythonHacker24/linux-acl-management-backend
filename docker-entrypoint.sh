#!/bin/bash

# Load env vars
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Run the binary
exec ./bin/laclm --config config.yaml
