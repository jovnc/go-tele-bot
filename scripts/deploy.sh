#!/bin/bash

# Deploy script for go-tele-bot
# Prerequisites: Secrets and IAM are managed via Terraform (./terraform)

set -e

# Check that gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo "Error: gcloud could not be found"
    exit 1
fi

# Check that project is set on gcloud
if ! gcloud config get project &> /dev/null; then
    echo "Error: project is not set on gcloud."
    exit 1
fi

# Load environment variables from .env
set -a
source .env
set +a

if [ -z "$WEBHOOK_URL" ]; then
    echo "Error: WEBHOOK_URL is not set"
    exit 1
fi

# Deploy to Google Cloud Run
# Secrets are managed by Terraform and referenced here
gcloud run deploy go-tele-bot \
    --source . \
    --region asia-southeast1 \
    --allow-unauthenticated \
    --set-secrets=BOT_TOKEN=bot-token:latest,GOOGLE_CREDENTIALS_JSON=google-credentials:latest,OPEN_AI_API_KEY=open-ai-api-key:latest \
    --set-env-vars="WEBHOOK_URL=$WEBHOOK_URL" \
    --set-env-vars="VALID_USERNAMES=$VALID_USERNAMES"
