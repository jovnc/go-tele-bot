# Go Telegram Bot - Google Drive Folder Sharing

A Telegram bot built with Go that allows you to share Google Drive folders to specified email addresses directly through Telegram.

## Features

- 📂 Select from pre-configured Google Drive folders
- ✅ Multi-select support with toggle buttons
- 📧 Share folders to any email address
- 🚀 Webhook-based for efficient message handling
- 🐳 Docker support for easy deployment

## Quick Start

### Prerequisites

- Go 1.23 or later
- A Telegram Bot Token (get one from [@BotFather](https://t.me/BotFather))
- A publicly accessible URL for webhooks (e.g., using ngrok for local development)

### 1. Clone and Install

```bash
git clone <your-repo-url>
cd go-tele-bot
go mod download && go mod verify
```

### 2. Configure Environment

Create a `.env` file in the project root with the following variables:

| Variable      | Description                                 | Required |
| ------------- | ------------------------------------------- | -------- |
| `BOT_TOKEN`   | Your Telegram bot token from BotFather      | Yes      |
| `WEBHOOK_URL` | Public URL where Telegram will send updates | Yes      |
| `PORT`        | Server port (default: 8080)                 | No       |

### 3. Configure Folders

Edit `data/folders.json` to add your Google Drive folders:

```json
[
  {
    "folder_name": "My Folder",
    "folder_id": "your_google_drive_folder_id"
  }
]
```

### 4. Run the Bot

```bash
make run
```

## Deploy to Google Cloud Run

1. Set project and enable APIs

```bash
gcloud config set project YOUR_PROJECT_ID
gcloud services enable run.googleapis.com secretmanager.googleapis.com
```

2. Create secrets and store in google cloud secrets

```bash
echo -n "YOUR_BOT_TOKEN" | gcloud secrets create bot-token --data-file=-
gcloud secrets create google-credentials --data-file=credentials.json
```

3. Deploy to google cloud run

```bash
gcloud run deploy go-tele-bot \
    --source . \
    --region asia-southeast1 \
    --allow-unauthenticated \
    --set-secrets=BOT_TOKEN=bot-token:latest,GOOGLE_CREDENTIALS_JSON=google-credentials:latest \
    --set-env-vars="WEBHOOK_URL=https://placeholder.run.app"
```

4. Update with actual cloud run URL

```bash
CLOUD_RUN_URL=$(gcloud run services describe go-tele-bot --region us-central1 --format='value(status.url)')
gcloud run services update go-tele-bot --region us-central1 --set-env-vars="WEBHOOK_URL=$CLOUD_RUN_URL"
```
