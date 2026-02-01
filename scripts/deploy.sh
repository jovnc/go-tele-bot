#!/bin/bash

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

# Load all variables from .env
set -a
source .env
set +a 

if [ -z "$BOT_TOKEN" ]; then
    echo "Error: BOT_TOKEN is not set"
    exit 1
fi

if [ -z "$OPEN_AI_API_KEY" ]; then
    echo "Error: OPEN_AI_API_KEY is not set"
    exit 1
fi

if [ -z "$WEBHOOK_URL" ]; then
    echo "Error: WEBHOOK_URL is not set"
    exit 1
fi

# Get project number
# PROJECT_NUMBER=$(gcloud projects describe "$(gcloud config get project)" --format="value(projectNumber)")

# Activate services (only needed once)
# gcloud services enable secretmanager.googleapis.com drive.googleapis.com run.googleapis.com cloudbuild.googleapis.com run.googleapis.com

# Create secrets (only needed once)
# echo -n $BOT_TOKEN | gcloud secrets create bot-token --data-file=-
# echo -n $OPEN_AI_API_KEY | gcloud secrets create open-ai-api-key --data-file=-
# gcloud secrets create google-credentials --data-file=credentials.json

# Set iam-permission for the service account (only needed once)
# gcloud secrets add-iam-policy-binding google-credentials \
#     --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
#     --role="roles/secretmanager.secretAccessor"

# gcloud secrets add-iam-policy-binding bot-token \
#     --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
#     --role="roles/secretmanager.secretAccessor"

# gcloud secrets add-iam-policy-binding open-ai-api-key \
#     --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
#     --role="roles/secretmanager.secretAccessor"

# Deploy the bot to Google Cloud Run
gcloud run deploy go-tele-bot \
    --source . \
    --region asia-southeast1 \
    --allow-unauthenticated \
    --set-secrets=BOT_TOKEN=bot-token:latest,GOOGLE_CREDENTIALS_JSON=google-credentials:latest,OPEN_AI_API_KEY=open-ai-api-key:latest \
    --set-env-vars="WEBHOOK_URL=$WEBHOOK_URL" \
    --set-env-vars="VALID_USERNAMES=$VALID_USERNAMES"
