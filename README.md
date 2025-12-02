# Go Telegram Bot - Google Drive Folder Sharing

A Telegram bot built with Go that allows you to share Google Drive folders to specified email addresses directly through Telegram.

## Features

- 📂 Select from pre-configured Google Drive folders
- 📧 Share folders to any email address
- 🚀 Webhook-based for efficient message handling
- 🐳 Docker support for easy deployment

## Quick Start (Local Development)

### Prerequisites

- Go 1.25 or later
- A Telegram Bot Token (get one from [@BotFather](https://t.me/BotFather))
- A publicly accessible URL for webhooks (e.g., ngrok for local development)

### 1. Clone and Install

```bash
git clone <your-repo-url>
cd go-tele-bot
go mod download && go mod verify
```

### 2. Configure Environment

Create a `.env` file in the project root by copying from `.env.example` and filling in the values.

### 3. Configure Folders

Edit `data/folders.json` by copying from `data/folders.json.example` and filling in the values.

```json
[
  {
    "folder_name": "My Folder",
    "folder_id": "your_google_drive_folder_id"
  }
]
```

> Note: You can get the folder ID by clicking on the folder in Google Drive and copying the ID from the URL.
> Example: `https://drive.google.com/drive/u/0/folders/1234567890`
> The folder ID is `1234567890`.

### 4. Run the Bot

```bash
make run
```

## Quick Start (Docker)

### Prerequisites

- Docker installed and running
- `.env` file with the values filled in
- `data/folders.json` file with the values filled in

### Build and Run Docker Image

This will build the Docker image and run the container.

```bash
make docker-run
```

## Quick Start (Google Cloud Run)

### Prerequisites

- `gcloud` CLI installed and authenticated (https://cloud.google.com/sdk/docs/install)
- A Telegram Bot Token (get one from [@BotFather](https://t.me/BotFather))
- `data/folders.json` file with the values filled in
- `credentials.json` file with service account credentials (create one from [Google Cloud Console](https://console.cloud.google.com/iam-admin/serviceaccounts))

> Note: service account needs to have the following permissions:
>
> - Google Drive API
> - Folder permissions to the folders to be shared (write/owner access)

### Deploy Setup for Google Cloud Run

This will setup the deployment for Google Cloud Run. Run this in the first deployment to google cloud run.

```bash
make deploy-setup
```

### Deploy to Google Cloud Run

This will update the deployment for Google Cloud Run. Run this to update the bot.

```bash
make deploy
```
