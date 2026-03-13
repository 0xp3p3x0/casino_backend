# Casino Backend

Simple Go backend using Gin, PostgreSQL, and JWT authentication.

## Features

- User registration and login
- JWT-based authentication
- Role-based access control (admin/user)
- WebSocket endpoint with token authentication
- PostgreSQL database managed by GORM

## Directory structure

```
cmd/api/main.go
internal/
  config/
  db/
  handlers/
  middleware/
  models/
  services/
  websocket/
  server/
```

## Getting Started

1. Copy `.env.example` to `.env` and fill in your database credentials.
2. Run `go mod tidy` to download dependencies.
3. Start the PostgreSQL database and ensure it is reachable.
4. Run the server:
   ```bash
   go run ./cmd/api
   ```
5. Use `POST /api/v1/auth/register` and `/api/v1/auth/login` to create accounts and obtain tokens.

## WebSocket

Open a WebSocket connection to `/ws?token=<JWT>` to receive profile and health updates.

## Notes

- Passwords are hashed with bcrypt.
- JWT tokens expire after 24 hours.
- Administrative routes are protected via `isAdmin=true` middleware.
