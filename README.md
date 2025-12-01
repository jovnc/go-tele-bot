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

- Go 1.25.4 or later
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

## Bot Commands

| Command  | Description                   |
| -------- | ----------------------------- |
| `/share` | Start the folder sharing flow |

## How It Works

1. User sends `/share` command
2. Bot displays available folders as selectable buttons
3. User toggles folders they want to share
4. User clicks "Done" when finished selecting
5. Bot prompts for an email address
6. User enters the email to share the folders with
