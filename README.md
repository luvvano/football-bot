# Football Bot ⚽

Telegram MiniApp for organizing football games in group chats.

## Features

- 📅 Create and manage football events
- 👥 Player registration with guest support
- ⚖️ Auto-balanced team formation based on skill ratings
- ⭐ Player rating system (speed, dribbling, shot)
- 📊 Statistics and leaderboards
- 🏟️ Venue management

## Architecture

```
football-bot/
├── backend/           # Go monolith (API + Bot)
│   ├── cmd/
│   │   ├── api/       # REST API server
│   │   └── bot/       # Telegram bot
│   └── internal/
│       ├── handlers/  # HTTP & bot handlers
│       ├── services/  # Business logic
│       ├── models/    # Data models
│       └── middleware/
├── miniapp/           # React MiniApp (TWA)
│   └── src/
│       ├── components/
│       ├── pages/
│       ├── hooks/
│       └── api/
├── migrations/        # PostgreSQL migrations
└── docker/            # Docker configs
```

## Tech Stack

- **Backend:** Go 1.22 + Chi router + telebot v3
- **Database:** PostgreSQL 16
- **Frontend:** React 18 + Vite + Telegram Web App SDK
- **Styling:** TailwindCSS

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- PostgreSQL 16
- Telegram Bot Token

### Development

```bash
# Backend
cd backend && go run ./cmd/bot

# MiniApp
cd miniapp && pnpm install && pnpm dev
```

### Environment Variables

```env
# Database
DATABASE_URL=postgres://user:pass@localhost:5432/football_bot

# Telegram
TELEGRAM_BOT_TOKEN=your_token
TELEGRAM_WEBAPP_URL=https://your-miniapp-url.com

# API
API_PORT=8080
```

## Bot Commands

| Command | Description |
|---------|-------------|
| `/start` | Register community, show MiniApp link |
| `/event` | Create new game |
| `/add` | Join the game |
| `/remove` | Leave the game |
| `/info` | Current game info |
| `/events` | List upcoming games |
| `/teams` | Form balanced teams |
| `/finish` | End game, trigger ratings |
| `/app` | Open MiniApp |

## License

MIT
