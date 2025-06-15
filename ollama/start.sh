#!/bin/sh

# This script starts the Ollama server in the background,
# waits for it to be ready, and then pulls the phi3:mini model.

# Curl seems to not be required but not included (?).
apt-get update && apt-get install -y curl

# Start Ollama server in background.
ollama serve &

# Wait for server to be up.
echo "Waiting for Ollama to start..."
until curl -s http://localhost:11434/ > /dev/null; do
  sleep 1
done

# Pull the model.
echo "Pulling phi3:mini model..."
ollama pull phi3:mini

# Mark the service as ready
touch /tmp/ready

# Keep the container running
wait
