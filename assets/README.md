# Chirpy Server

A RESTful API server built in Go that provides a Twitter-like platform for posting short messages (chirps).

## Overview

Chirpy Server is a full-featured backend application that handles user authentication, content management, and premium user features. It provides a solid foundation for building a social media platform with secure access control and data persistence.

## Features

- **User Authentication**
  - Register new users
  - Login with email/password
  - JWT-based authentication
  - Refresh token support
  - Token revocation

- **Content Management**
  - Create chirps (short messages up to 140 characters)
  - List all chirps with sorting options
  - Filter chirps by author
  - Get a specific chirp by ID
  - Delete chirps (only by the owner)
  - Profanity filtering

- **Premium Features**
  - Chirpy Red subscription support via webhooks
  - Polka API integration for payments

- **Administration Tools**
  - Reset database in development mode
  - Server metrics and monitoring

## Technical Stack

- **Language**: Go
- **Database**: PostgreSQL
- **ORM**: SQLC for type-safe SQL queries
- **Authentication**: JWT (JSON Web Tokens)
- **API**: RESTful HTTP endpoints
- **Configuration**: Environment variables with godotenv
- **Migration**: Goose for database schema management

## API Endpoints

### Authentication
- `POST /api/users` - Register a new user
- `POST /api/login` - Login and get access/refresh tokens
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke a refresh token

### User Management
- `PUT /api/users` - Update user email and password

### Chirps
- `POST /api/chirps` - Create a new chirp
- `GET /api/chirps` - List all chirps (with optional sorting and filtering)
- `GET /api/chirps/{chirpID}` - Get a specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete a chirp (owner only)

### Admin
- `POST /admin/reset` - Reset the database (dev mode only)
- `GET /admin/metrics` - View server metrics

### Webhooks
- `POST /api/polka/webhooks` - Handle Polka payment webhooks

### System
- `GET /api/healthz` - Health check endpoint

## Project Structure

```
.
├── assets/
│   └── logo.png
├── internal/
│   ├── auth/
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── database/
│   │   ├── db.go
│   │   ├── models.go
│   │   └── 001_users.sql.go
│   └── util/
│       └── string_utils.go
├── sql/
│   ├── queries/
│   │   └── 001_users.sql
│   └── schema/
│       └── 001_users.sql
├── .env
├── .gitignore
├── chirps.go
├── go.mod
├── go.sum
├── index.html
├── json.go
├── main.go
├── metrics.go
├── polka.go
├── readiness.go
├── reset.go
├── sqlc.yaml
└── users.go
```

## Getting Started

### Prerequisites
- Go 1.23.5 or later
- PostgreSQL database
- Goose (for migrations)
- SQLC (for generating database code)

### Environment Setup
Create a .env file with:
```
DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
PLATFORM=dev
SECRET="your_jwt_secret_key"
POLKA_KEY="your_polka_api_key"
```

### Database Setup
1. Create a PostgreSQL database
2. Run migrations:
```bash
goose -dir sql/schema postgres "your_connection_string" up
```

### Run the Server
```bash
go run .
```

The server will start on port 8080.

## Testing
Run the test suite:
```bash
go test ./...
```

## License
MIT License