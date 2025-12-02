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

# Deploy the bot to Google Cloud Run
gcloud run deploy go-tele-bot \
    --source . \
    --region asia-southeast1

# Update the bot
# gcloud run services update go-tele-bot \
#     --region asia-southeast1 \
#     --set-env-vars "VALID_USERNAMES=test0;test1" \
#     --set-env-vars="WEBHOOK_URL=https://go-tele-bot-XXX.asia-southeast1.run.app"
