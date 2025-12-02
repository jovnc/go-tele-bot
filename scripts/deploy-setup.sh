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

# Activate services
gcloud services enable secretmanager.googleapis.com drive.googleapis.com run.googleapis.com cloudbuild.googleapis.com run.googleapis.com

# Create secrets
echo -n $BOT_TOKEN | gcloud secrets create bot-token --data-file=-
gcloud secrets create google-credentials --data-file=credentials.json

# Set iam-permission for the service account
gcloud secrets add-iam-policy-binding google-credentials \
    --member="serviceAccount:service-${PROJECT_NUMBER}@developer.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding bot-token \
    --member="serviceAccount:service-${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

# Deploy the bot to Google Cloud Run
gcloud run deploy go-tele-bot \
    --source . \
    --region asia-southeast1 \
    --allow-unauthenticated \
    --set-secrets=BOT_TOKEN=bot-token:latest,GOOGLE_CREDENTIALS_JSON=google-credentials:latest \
    --set-env-vars="WEBHOOK_URL=https://placeholder.run.app"
