# Content Moderation System

A Go-based content moderation service that uses AI-powered analysis to automatically moderate text and image content.

## Features

- **Content Upload**: Upload text and image content for moderation
- **AI-Powered Moderation**: Uses Google Gemini for intelligent content analysis
- **Async Processing**: Background job processing with Redis-backed queue (Asynq)
- **Image Storage**: ImageKit integration for image uploads and management
- **Admin Review**: Manual review workflow for moderated content

## Tech Stack

- **Go** - Core language
- **Gorilla Mux** - HTTP router
- **GORM** - ORM for PostgreSQL
- **Asynq** - Redis-based async task queue
- **Google Gemini** - AI content moderation
- **ImageKit** - Image upload and CDN
- **PostgreSQL** (Neon) - Primary database
- **Redis** (Upstash) - Queue backend

## Prerequisites

- Go 1.25+
- PostgreSQL database (or Neon account)
- Redis instance (or Upstash account)
- ImageKit account
- Google Gemini API key

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd moderation_go
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Edit `.env` with your actual credentials.

### 3. Install dependencies

```bash
go mod download
```

### 4. Run the application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`.

## Project Structure

```
moderation_go/
├── api/
│   ├── handlers/     # HTTP request handlers
│   └── routes/       # Route definitions
├── cmd/
│   └── main.go       # Application entry point
├── internal/
│   ├── database/     # Database connection
│   ├── models/       # Data models
│   └── queue/        # Async job processing
│       ├── workers/  # Job handlers (text, image, aggregation)
│       └── worker-client/
├── utils/
│   ├── cors/         # CORS middleware
│   ├── gemini/       # Gemini AI client
│   └── imagekit/     # ImageKit client
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/upload/content` | Upload content for moderation |
| GET | `/content` | Get all content |
| PATCH | `/content/update` | Update content status (admin) |

## License

MIT
