# Message Provider Go

A robust HTTP service written in Go that manages and processes messages with scheduled background jobs.

## Features

- RESTful API for message operations
- PostgreSQL database integration with GORM
- Background job scheduler for automated message processing
- Transaction management for data consistency
- Scheduler control endpoints for runtime job management
- Docker containerization with Docker Compose
- Gin web framework for high performance
- CORS support
- Database connection middleware
- Graceful shutdown with signal handling
- Health check endpoints

## API Endpoints

### Health Check
- `GET /health` - Returns server health status
- `GET /api/v1/health` - Returns detailed health status with timestamp

### Messages
- `GET /api/v1/message/` - Retrieve 2 sent messages using transaction

### Scheduler Management
- `POST /api/v1/scheduler/job/start` - Start the FetchMessagesJob
- `POST /api/v1/scheduler/job/stop` - Stop the FetchMessagesJob
- `POST /api/v1/scheduler/job/restart` - Restart the FetchMessagesJob
- `GET /api/v1/scheduler/job/status` - Get status of all scheduled jobs
- `POST /api/v1/scheduler/stop` - Stop the entire scheduler

## Quick Start

### Using Docker Compose (Recommended)

1. Start the application with PostgreSQL:
```bash
docker compose up -d
```

2. Check the logs:
```bash
docker compose logs -f backend
```

3. Stop the services:
```bash
docker compose down
```

### Manual Setup

1. Install dependencies:
```bash
go mod download
```

2. Set up PostgreSQL database and run the init.sql script

3. Configure environment variables (see Configuration section)

4. Run the server:
```bash
go run main.go
```

The server will start on port 8080 by default.

## Configuration

The application can be configured using environment variables:

### Server Configuration
- `PORT` - Server port (default: 8080)
- `ENV` - Environment (default: development)
- `LOG_LEVEL` - Log level: debug, info, warn, error (default: info)

### Database Configuration
- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_NAME` - Database name (default: provider)
- `DB_USER` - Database user (default: admin)
- `DB_PASSWORD` - Database password (default: Aa123456)
- `DB_SSL_MODE` - SSL mode (default: disable)

### Scheduler Configuration
The FetchMessagesJob runs every 10 seconds by default and processes 2 unsent messages per execution.

## Example Usage

### Get Messages
```bash
curl http://localhost:8080/api/v1/message/
```

### Scheduler Control

#### Start the message processing job
```bash
curl -X POST http://localhost:8080/api/v1/scheduler/job/start
```

#### Stop the message processing job
```bash
curl -X POST http://localhost:8080/api/v1/scheduler/job/stop
```

#### Restart the message processing job
```bash
curl -X POST http://localhost:8080/api/v1/scheduler/job/restart
```

#### Get scheduler job status
```bash
curl http://localhost:8080/api/v1/scheduler/job/status
```

Response example:
```json
{
  "status": "success",
  "jobs": [
    {
      "name": "fetch-messages",
      "interval": "10s"
    }
  ],
  "count": 1
}
```

#### Stop the entire scheduler
```bash
curl -X POST http://localhost:8080/api/v1/scheduler/stop
```

### Health Check
```bash
curl http://localhost:8080/health
# or
curl http://localhost:8080/api/v1/health
```

## Project Structure

```
├── main.go                     # Application entry point
├── go.mod                      # Go module file
├── go.sum                      # Go module checksums
├── Dockerfile                  # Docker image definition
├── docker-compose.yaml         # Docker Compose configuration
├── init.sql                    # Database initialization script
├── Makefile                    # Build automation
├── build/
│   └── goapp                   # Compiled binary
└── internal/
    ├── config/
    │   └── config.go          # Configuration management
    ├── database/
    │   └── connection.go      # Database connection setup
    ├── handlers/
    │   ├── message.go         # Message HTTP handlers
    │   └── scheduler.go       # Scheduler HTTP handlers
    ├── middleware/
    │   └── middleware.go      # HTTP middleware (CORS, DB)
    ├── models/
    │   └── message.go         # Message data model
    ├── repository/
    │   └── message.go         # Database repository layer
    ├── scheduler/
    │   ├── scheduler.go       # Job scheduler implementation
    │   └── jobs.go            # Scheduled job definitions
    ├── server/
    │   └── server.go          # Server setup and routes
    └── service/
        └── message.go         # Business logic layer
```

## Architecture

### Layered Architecture
The application follows a clean layered architecture:

1. **Handler Layer** (`internal/handlers/`) - HTTP request/response handling
2. **Service Layer** (`internal/service/`) - Business logic and transaction management
3. **Repository Layer** (`internal/repository/`) - Database operations
4. **Model Layer** (`internal/models/`) - Data structures

### Background Jobs
The scheduler runs background jobs for automated message processing:
- **FetchMessagesJob** - Automatically fetches unsent messages, processes them, and updates their status
- Configurable job intervals
- Runtime control via REST API

### Database
- PostgreSQL 15 with GORM ORM
- Transaction support for data consistency
- Connection pooling and middleware

## Dependencies

- [gin-gonic/gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- [gorm.io/gorm](https://gorm.io/) - ORM library for Go
- [gorm.io/driver/postgres](https://gorm.io/driver/postgres) - PostgreSQL driver for GORM
- [go-playground/validator](https://github.com/go-playground/validator) - Struct validation
- [lib/pq](https://github.com/lib/pq) - PostgreSQL driver

## Development

### Run in development mode
```bash
ENV=development LOG_LEVEL=debug go run main.go
```

### Build for production
```bash
go build -o build/goapp main.go
```

### Using Makefile
```bash
# Install dependencies
make deps

# Build the application
make build

# Run tests
make test

# Clean build artifacts
make clean
```

### Docker Commands
```bash
# Build and start all services
docker compose up -d

# View logs
docker compose logs -f backend
docker compose logs -f postgres

# Restart a specific service
docker compose restart backend

# Stop all services
docker compose down

# Rebuild and start
docker compose up -d --build
```

## How It Works

1. **Application Startup**:
   - Loads configuration from environment variables
   - Establishes PostgreSQL database connection
   - Initializes scheduler with scheduled jobs
   - Starts HTTP server on configured port

2. **Message Processing**:
   - The scheduler runs `FetchMessagesJob` every 10 seconds
   - Job fetches 2 unsent messages from the database
   - Processes each message (simulates sending)
   - Updates message status to "sent" in the database
   - Uses transactions to ensure data consistency

3. **Runtime Control**:
   - Scheduler can be controlled via REST API endpoints
   - Start/stop/restart individual jobs or entire scheduler
   - Monitor job status and intervals in real-time

4. **Graceful Shutdown**:
   - Handles SIGINT and SIGTERM signals
   - Stops scheduler and completes running jobs
   - Closes database connections
   - Shuts down HTTP server gracefully

## License

MIT