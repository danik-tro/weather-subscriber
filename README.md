# Weather Subscription Service
[![codecov](https://codecov.io/github/danik-tro/weather-subscriber/branch/master/graph/badge.svg?token=VHrmK2NoIP)](https://codecov.io/github/danik-tro/weather-subscriber)

A Go-based service that allows users to subscribe to weather updates for their preferred cities. The service sends weather updates via email at configurable intervals (hourly or daily).

## Approach
Due to limited time, it was decided to use a cron job and a custom event publisher abstraction instead of implementing RabbitMQ or other queues and background task managers at this stage. The cron job is triggered either hourly or daily to send weather updates to users.

To avoid hitting the external API unnecessarily, a caching mechanism was introduced. If multiple users subscribe to the same city, the weather data is fetched from the cache instead of making repeated API calls. This significantly optimizes performance.

The confirmation email for the subscription is sent asynchronously. In the current implementation, a simple custom publisher-subscriber model with worker routines was built. In a production environment, this setup could be extended with dedicated brokers and workers to handle larger volumes of data more efficiently.

For the backend architecture, the onion architecture pattern was chosen, utilizing use cases instead of service classes. This approach enabled the creation of single-purpose functions that perform exactly one operation at the handler level. Everything is abstracted to facilitate testing and to allow for easy replacement of implementations.

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

### Local Development
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

### Docker Environment
For Docker deployment, create a `.docker.env` file with the following configuration:

```env
# Database configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=weather_app
DB_SSL_MODE=disable
DB_AUTO_MIGRATE=true

# Redis configuration
REDIS_ADDRESS=redis:6379
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

# Base URL (use the service name from docker-compose)
BASE_URL=http://app:8080

# Application configuration
APP_HOST=0.0.0.0
APP_PORT=8080
```

Key differences in Docker environment:
- Database host is set to `postgres` (service name in docker-compose)
- Redis host is set to `redis` (service name in docker-compose)
- Base URL uses the service name `app` instead of localhost
- Application host is set to `0.0.0.0` to allow external connections
- Database and Redis credentials are simplified for development

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
POST /api/confirm/{confirmation_token}
```

### Unsubscribe
```http
POST /api/unsubscribe/{unsubscribe_token}
```

HTML pages
### Confirm Subscription
```http
GET /confirm/{confirmation_token}
```

### Unsubscribe
```http
GET /unsubscribe/{unsubscribe_token}
```

## Running the Service

### Local Development
1. Clone the repository:
```bash
git clone https://github.com/danik-tro/weather-subscriber.git
cd weather-subscriber
```

2. Install dependencies:
```bash
go mod download
```

3. Run the service:
```bash
go run main.go
```

### Docker Deployment
The project includes Docker Compose configuration for easy deployment:

1. Create the `.docker.env` file as described above
2. Build and start the services:
```bash
docker-compose up -d
```

This will start:
- The main application
- PostgreSQL database
- Redis cache

To view logs:
```bash
docker-compose logs -f
```

To stop the services:
```bash
docker-compose down
```

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
├── .env                # Environment configuration for local development
├── .docker.env         # Environment configuration for Docker
├── compose.yaml        # Docker Compose configuration
└── README.md          # This file
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
