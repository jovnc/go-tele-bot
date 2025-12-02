#!/bin/bash
# Script to run the Docker container
# Usage: ./docker-run.sh
#
# Prerequisites:
# - .env file with the values filled in
# - credentials.json file with the values filled in
# - data/folders.json file with the values filled in
# - docker installed and running

# Check if the .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found"
    exit 1
fi

# Check if the credentials.json file exists
if [ ! -f credentials.json ]; then
    echo "Error: credentials.json file not found"
    exit 1
fi

# Check if the data/folders.json file exists
if [ ! -f data/folders.json ]; then
    echo "Error: data/folders.json file not found"
    exit 1
fi

# Check if docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: docker could not be found"
    exit 1
fi

# Check if docker is running
if ! docker info &> /dev/null; then
    echo "Error: docker is not running"
    exit 1
fi

# Build the Docker image
docker build -t go-tele-bot .
if [ $? -ne 0 ]; then
    echo "Error: failed to build the Docker image"
    exit 1
fi

# Run the Docker container
docker run --rm \
    --env-file .env \
    -e GOOGLE_CREDENTIALS_JSON="$(cat credentials.json)" \
    -p 8080:8080 \
    go-tele-bot
if [ $? -ne 0 ]; then
    echo "Error: failed to run the Docker container"
    exit 1
fi
