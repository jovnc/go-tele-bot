# Go Telegram Bot - My Personal Assistant

A Telegram bot built with Go that provides:
- 📂 Google Drive folder sharing directly through Telegram
- 💻 Interactive LeetCode practice with AI assistance

## Features

### Google Drive Integration
- 📂 Select from pre-configured Google Drive folders
- 📧 Share folders to any email address
- 🔐 Secure service account authentication

### LeetCode Practice Mode
- 🎯 Generate random LeetCode problems or request specific topics
- 🤖 AI-powered guidance and hints (no direct solutions)
- 💬 Interactive problem-solving with conversation history
- 📊 Code analysis and complexity feedback

### Technical Features
- 🚀 Webhook mode for production (Cloud Run)
- 🔄 Polling mode for local development
- 🐳 Docker support for easy deployment
- 🔒 Username-based access control

## Quick Start (Local Development)

### Prerequisites

- Go 1.25 or later
- A Telegram Bot Token (get one from [@BotFather](https://t.me/BotFather))
- OpenAI API Key (for LeetCode features)
- Google Service Account credentials (for Drive sharing features)

### 1. Clone and Install

```bash
git clone <your-repo-url>
cd go-tele-bot
go mod download && go mod verify
```

### 2. Configure Environment

Create a `.env` file in the project root:

```env
BOT_TOKEN=your_telegram_bot_token
WEBHOOK_URL=https://your-webhook-url.com
PORT=8080
VALID_USERNAMES=username1;username2
OPEN_AI_API_KEY=your_openai_api_key
DEV=true
```

**Note:** Set `DEV=true` for local development (uses polling mode). Remove or set to `false` for production (uses webhook mode).

### 3. Configure Data Files

**Folders Configuration** - Edit `data/folders.json`:

```json
[
  {
    "folder_name": "My Folder",
    "folder_id": "your_google_drive_folder_id"
  }
]
```

**LeetCode System Prompt** - Optionally customize `data/lc_system_prompt.md` to change the AI assistant's behavior.

### 4. Add Google Credentials

Place your Google service account credentials in `credentials.json` at the project root.

### 5. Run the Bot

```bash
DEV=true && make run
```

## Quick Start (Docker)

### Prerequisites

- Docker installed and running
- `.env` file with `DEV=true` for local testing (optional)
- `credentials.json` file with Google service account credentials
- `data/folders.json` file with folder configuration

### Build and Run Docker Image

```bash
make docker-run
```

The bot will automatically use polling mode for local development when `DEV=true` is set in `.env`.

## Bot Commands

### General Commands
- `/start` - Start the bot and see available commands
- `/help` - Show help information

### LeetCode Commands
- `/lc` - Get a random LeetCode problem
- `/lc <topic>` - Get a LeetCode problem about a specific topic (e.g., `/lc arrays`, `/lc dynamic programming`)
- `/exit` - Exit LeetCode practice mode

**Examples:**
```
/lc
/lc binary trees
/lc dynamic programming
/lc two pointers
```

Once in LC mode, simply send messages to interact with the AI assistant. The assistant will:
- Guide you through the problem
- Provide hints when you're stuck
- Analyze your solution's time/space complexity
- Suggest optimizations

### Google Drive Commands
- `/share` - Start the folder sharing flow
  1. Select a folder from the list
  2. Enter the email address to share with
  3. Confirm the sharing

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `BOT_TOKEN` | Telegram bot token from BotFather | Yes | - |
| `WEBHOOK_URL` | Public URL for webhooks (production only) | Yes (prod) | - |
| `PORT` | Server port | No | 8080 |
| `VALID_USERNAMES` | Semicolon-separated list of allowed usernames | Yes | - |
| `OPEN_AI_API_KEY` | OpenAI API key for LeetCode features | Yes | - |
| `DEV` | Set to `true` for local development (polling mode) | No | false |

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

### Deploy to Google Cloud Run

This will deploy telegram bot for Google Cloud Run. Run this to start/update the bot.

```bash
make deploy
```
