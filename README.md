# Weather Subscription Service

A Go-based service that allows users to subscribe to weather updates for their preferred cities. The service sends weather updates via email at configurable intervals (hourly or daily).

## Features

- **Weather Subscriptions**
  - Subscribe to weather updates for any city
  - Choose between hourly or daily updates
  - Email confirmation required for new subscriptions
  - Easy unsubscribe option

- **Weather Updates**
  - Current temperature
  - Weather conditions
  - Humidity levels
  - Beautiful HTML email templates

- **Scheduling**
  - Daily updates sent at noon (12:00)
  - Hourly updates sent at the start of each hour
  - Configurable update frequencies

- **Infrastructure**
  - PostgreSQL database for subscription storage
  - Redis for weather data caching
  - SMTP email service integration
  - RESTful API endpoints

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 15
- Redis 5
- SMTP server access

## Configuration

Create a `.env` file in the project root with the following variables:

```env
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=weather_app
DB_SSL_MODE=disable
DB_AUTO_MIGRATE=true

# Redis configuration
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Weather API configuration
WEATHER_API_KEY=your_api_key

# SMTP configuration
SMTP_HOST=your_smtp_host
SMTP_PORT=587
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password
SMTP_FROM=your_from_email

# Base URL
BASE_URL=http://localhost:8080
```

## API Endpoints

### Subscribe to Weather Updates
```http
POST /api/subscribe
Content-Type: application/json

{
    "email": "user@example.com",
    "city": "London",
    "frequency": "daily"  // or "hourly"
}
```

### Get Current Weather
```http
GET /api/weather?city=London
```

### Confirm Subscription
```http
GET /confirm/{confirmation_token}
```

### Unsubscribe
```http
GET /unsubscribe/{unsubscribe_token}
```

## Running the Service

1. Clone the repository:
```bash
git clone https://github.com/yourusername/weather-subscriber.git
cd weather-subscriber
```

2. Install dependencies:
```bash
go mod download
```

3. Run the service:
```bash
go run cmd/main.go
```

## Docker Support

The project includes Docker Compose configuration for easy deployment:

```bash
docker-compose up -d
```

This will start:
- The main application
- PostgreSQL database
- Redis cache

## Project Structure

```
.
├── cmd/
│   └── main.go           # Application entry point
├── pkg/
│   ├── domain/          # Domain models and interfaces
│   ├── infrastructure/  # Infrastructure implementations
│   │   ├── background_job/  # Background job service
│   │   ├── db/          # Database implementations
│   │   ├── email_service/   # Email service
│   │   └── events/      # Event handling
│   ├── external/        # External service integrations
│   │   └── weather/     # Weather API client
│   └── presenter/       # API handlers and routes
├── templates/           # Email templates
├── .env                # Environment configuration
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.