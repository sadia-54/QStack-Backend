# QStack Backend

A professional Q&A platform backend (similar to Stack Overflow) built with Go, featuring user authentication, question & answer management, voting, personalized feeds, and async email processing.

[![Go Version](https://img.shields.io/badge/go-1.25.5-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![PostgreSQL](https://img.shields.io/badge/postgresql-17+-336791?logo=postgresql)](https://www.postgresql.org/)
[![RabbitMQ](https://img.shields.io/badge/rabbitmq-3.12+-FF6600?logo=rabbitmq)](https://www.rabbitmq.com/)

---

## Features

- **User Management** — Registration, JWT authentication, email verification, password reset
- **Q&A System** — Create questions, answers, and comments with full CRUD operations
- **Voting System** — Upvote/downvote questions with toggle behavior
- **Tag System** — Auto-created tags, tag-based search, and filtering
- **Personalized Feeds** — AI-driven question recommendations based on user tag interests
- **Activity Tracking** — User profiles with activity history and community statistics
- **Async Email** — Non-blocking email delivery via RabbitMQ and Mailpit
- **File Uploads** — Profile image upload support

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP Layer                           │
│                    Echo Framework (Go)                       │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                     Business Layer                          │
│                    Services (Logic)                         │
└──┬──────────────────────────────┬──────────────────────────┘
   │                              │
   │         ┌────────────────────▼────────────────────┐      │
   │         │          Message Queue (RabbitMQ)       │      │
   │         │     Producer → Queue → Consumer         │      │
   │         └────────────────────┬────────────────────┘      │
   │                              │                           │
┌──▼──────────────────┐  ┌────────▼──────────────────────┐   │
│   Data Layer        │  │        Email Worker           │   │
│ (GORM Repositories) │  │      (SMTP via Mailpit)       │   │
└──┬──────────────────┘  └───────────────────────────────┘   │
   │                                                          │
┌──▼──────────────────────────────────────────────────────────┐
│                     PostgreSQL Database                     │
└─────────────────────────────────────────────────────────────┘
```

**Layered Design:**

- **Handlers** → HTTP request/response (Echo)
- **Services** → Business logic and validation
- **Repositories** → Data access (GORM ORM)
- **Models** → Domain entities and DTOs
- **Queue** → Async job processing (RabbitMQ)

---

## Quick Start

### Prerequisites

| Requirement                                   | Description                                                |
| --------------------------------------------- | ---------------------------------------------------------- |
| [Go](https://go.dev/) 1.25.5                                     |
| [PostgreSQL](https://www.postgresql.org/) 17 |
| [RabbitMQ](https://www.rabbitmq.com/)         | Message queue (Docker image: `rabbitmq:3-management`)      |
| [Mailpit](https://github.com/axllent/mailpit) | SMTP testing server (Docker image: `axllent/mailpit`)      |
| [Docker](https://www.docker.com/) (optional)  | Used to run RabbitMQ and Mailpit via `docker-compose`      |

### 1. Clone the Repository

```bash
git clone https://github.com/sadia-54/QStack-Backend.git
cd QStack-Backend
```

### 2. Configure Environment

Create a `.env` file in the project root:

```env
# Application
APP_PORT=8080
APP_BASE_URL=http://localhost:3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=qstack
DB_SSLMODE=disable

# Authentication
JWT_SECRET=your-super-secret-jwt-key-change-this

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Mailpit (SMTP)
MAILPIT_HOST=localhost
MAILPIT_PORT=1025
```

### 3. Start Infrastructure (Docker)

Start RabbitMQ and Mailpit using Docker Compose:

```bash
docker-compose -f docker/docker-compose.yml up -d
```

This starts:

- **RabbitMQ** — `amqp://localhost:5672` (Management UI: `http://localhost:15672`)
- **Mailpit** — SMTP on `localhost:1025` (Web UI: `http://localhost:8025`)

### 4. Run Database Migrations

```bash
go run cmd/migrator/main.go -action up
```

### 5. Start the Server

```bash
# API Server (port 8080)
go mod tidy
go run cmd/server/main.go

# Email Worker (background process)
go run cmd/worker/main.go
```

### 6. Verify Installation

```bash
curl http://localhost:8080/health
# Expected: {"status":"ok","database":"connected"}
```

---

## Project Structure

```
QStack-Backend/
├── cmd/                          # Application entry points
│   ├── server/main.go            # HTTP API server (Echo)
│   ├── worker/main.go            # Background email worker
│   └── migrator/main.go          # Database migration CLI
│
├── internal/
│   ├── api/
│   │   ├── handlers/             # HTTP request handlers
│   │   │   ├── auth.go           # Authentication endpoints
│   │   │   ├── question.go       # Question CRUD + voting
│   │   │   ├── answer.go         # Answer CRUD + accept
│   │   │   ├── comment.go        # Comment CRUD
│   │   │   ├── user.go           # User profile + activity
│   │   │   └── upload.go         # File upload
│   │   ├── routes/               # Route registration
│   │   └── middleware/
│   │       └── jwt.go            # JWT authentication middleware
│   │
│   ├── config/                   # Configuration & database
│   │   ├── config.go             # Environment loading
│   │   └── db.go                 # GORM connection
│   │
│   ├── models/
│   │   ├── domains/              # Database entities
│   │   └── dtos/                 # Request/Response DTOs
│   │
│   ├── repositories/             # Data access layer (GORM)
│   ├── services/                 # Business logic
│   ├── queue/                    # RabbitMQ producer/consumer
│   ├── workers/                  # Background job processors
│   └── validator/                # Request validation
│
├── migrations/                   # SQL migration files
├── docker/
│   └── docker-compose.yml        # Infrastructure services
├── uploads/                      # User uploaded files
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── DOCUMENTATION.md              # Complete API documentation
└── README.md                     # This file
```

---

## 🔌 API Endpoints

**Base URL:** `http://localhost:8080/api/v1`

### Health Check

| Method | Endpoint  | Description                  |
| ------ | --------- | ---------------------------- |
| `GET`  | `/health` | API + Database health status |

### Authentication

| Method | Endpoint                       | Description                     |
| ------ | ------------------------------ | ------------------------------- |
| `POST` | `/auth/signup`                 | Register new user               |
| `POST` | `/auth/login`                  | Authenticate user               |
| `POST` | `/auth/logout`                 | Clear auth cookies              |
| `GET`  | `/auth/verify-email?token=xxx` | Verify email address            |
| `POST` | `/auth/forgot-password`        | Request password reset          |
| `POST` | `/auth/reset-password`         | Reset password with token       |
| `POST` | `/auth/change-password`        | Change password (auth required) |

### Questions

| Method   | Endpoint              | Description          | Auth        |
| -------- | --------------------- | -------------------- | ----------- |
| `GET`    | `/questions`          | Public question feed | No          |
| `GET`    | `/questions/:id`      | Get single question  | No          |
| `GET`    | `/questions/my-feed`  | Personalized feed    | Yes         |
| `GET`    | `/questions/my`       | User's questions     | Yes         |
| `POST`   | `/questions`          | Create question      | Yes         |
| `PUT`    | `/questions/:id`      | Update question      | Yes (owner) |
| `DELETE` | `/questions/:id`      | Delete question      | Yes (owner) |
| `POST`   | `/questions/:id/vote` | Vote on question     | Yes         |

### Answers

| Method   | Endpoint                         | Description   | Auth          |
| -------- | -------------------------------- | ------------- | ------------- |
| `GET`    | `/answers/question/:question_id` | Get answers   | No            |
| `POST`   | `/answers/question/:question_id` | Create answer | Yes           |
| `PUT`    | `/answers/:id`                   | Update answer | Yes (owner)   |
| `DELETE` | `/answers/:id`                   | Delete answer | Yes (owner)   |
| `PUT`    | `/answers/:id/accept`            | Accept answer | Yes (Q owner) |

### Comments

| Method   | Endpoint                      | Description    | Auth        |
| -------- | ----------------------------- | -------------- | ----------- |
| `GET`    | `/comments/answer/:answer_id` | Get comments   | No          |
| `POST`   | `/comments/answer/:answer_id` | Create comment | Yes         |
| `PUT`    | `/comments/:id`               | Update comment | Yes (owner) |
| `DELETE` | `/comments/:id`               | Delete comment | Yes (owner) |

### Users

| Method | Endpoint                 | Description          | Auth      |
| ------ | ------------------------ | -------------------- | --------- |
| `GET`  | `/users`                 | List users           | No        |
| `GET`  | `/users/:id/profile`     | User profile         | No        |
| `GET`  | `/users/me`              | Current user profile | Yes       |
| `PUT`  | `/users/profile`         | Update bio           | Yes       |
| `POST` | `/users/profile/image`   | Upload profile image | Yes       |
| `GET`  | `/users/:id/activity`    | User activity        | Yes (own) |
| `GET`  | `/users/community/stats` | Community statistics | No        |

### Tags & Upload

| Method | Endpoint        | Description          | Auth |
| ------ | --------------- | -------------------- | ---- |
| `GET`  | `/tags/popular` | Top 10 tags by usage | No   |
| `POST` | `/upload`       | Upload image file    | No   |

> **Full API Documentation:** See [DOCUMENTATION.md](DOCUMENTATION.md) for detailed request/response schemas and examples.

---

## Commands

### Development

```bash
# Run API server
go run cmd/server/main.go

# Run email worker
go run cmd/worker/main.go

# Run migrations
go run cmd/migrator/main.go -action up            # Apply all migrations
go run cmd/migrator/main.go -action down          # Rollback last migration
go run cmd/migrator/main.go -action version       # Show current version
go run cmd/migrator/main.go -action force -forceVersion N  # Force version
```

### Build

```bash
# Build binaries
go build -o bin/server cmd/server/main.go
go build -o bin/worker cmd/worker/main.go
go build -o bin/migrator cmd/migrator/main.go
```

### Docker (Infrastructure)

```bash
# Start RabbitMQ + Mailpit
docker-compose -f docker/docker-compose.yml up -d

# Stop services
docker-compose -f docker/docker-compose.yml down
```

---

## Configuration

All configuration is managed through environment variables or a `.env` file.

| Variable       | Default                              | Description              |
| -------------- | ------------------------------------ | ------------------------ |
| `APP_PORT`     | `8080`                               | HTTP server port         |
| `APP_BASE_URL` | `http://localhost:8080`              | Base URL for email links |
| `DB_HOST`      | `localhost`                          | PostgreSQL host          |
| `DB_PORT`      | `5432`                               | PostgreSQL port          |
| `DB_USER`      | `postgres`                           | Database user            |
| `DB_PASSWORD`  | _(required)_                         | Database password        |
| `DB_NAME`      | `qstack`                             | Database name            |
| `DB_SSLMODE`   | `disable`                            | PostgreSQL SSL mode      |
| `JWT_SECRET`   | _(required)_                         | JWT signing secret       |
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:5672/` | RabbitMQ connection URL  |
| `MAILPIT_HOST` | `localhost`                          | Mailpit SMTP host        |
| `MAILPIT_PORT` | `1025`                               | Mailpit SMTP port        |

---

## Database Schema

**10 Tables:**

- `users` — User accounts with email verification
- `tags` — Tag names (auto-created)
- `user_preferred_tags` — User tag preferences
- `questions` — Questions with voting
- `question_tags` — Question-tag mapping
- `answers` — Answers with acceptance status
- `comments` — Comments on answers
- `question_votes` — User votes (+1/-1)
- `password_reset_tokens` — Password reset tokens
- `email_verification_tokens` — Email verification tokens

> **Full Schema Details:** See [migrations/](migrations/) directory for SQL definitions.

---

## Security Features

- **JWT Authentication** — HS256 signed tokens with HTTP-only cookies
- **Password Hashing** — bcrypt with default cost factor
- **Token Hashing** — SHA256 hashing for tokens before database storage
- **Input Validation** — All requests validated with go-playground/validator
- **Generic Errors** — Do not reveal user existence on login/forgot-password
- **Owner Authorization** — Users can only modify their own content
- **Database Constraints** — Unique email/username, vote uniqueness, single accepted answer

---

## Tech Stack

| Layer                | Technology                   |
| -------------------- | ---------------------------- |
| **Language**         | Go 1.25.5                    |
| **Web Framework**    | Echo v4                      |
| **Database**         | PostgreSQL (via GORM)        |
| **ORM**              | GORM v1.31.1                 |
| **Migrations**       | golang-migrate/migrate v4    |
| **Message Queue**    | RabbitMQ (amqp091-go)        |
| **Email Testing**    | Mailpit                      |
| **Authentication**   | JWT (golang-jwt)             |
| **Validation**       | go-playground/validator v10  |
| **Password Hashing** | bcrypt (golang.org/x/crypto) |

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Additional Resources

- [ Full API Documentation](DOCUMENTATION.md) — Detailed request/response schemas
- [ Database Migrations](migrations/) — SQL schema definitions
- [ Docker Setup](docker/docker-compose.yml) — Infrastructure services

---

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

<div align="center">
  <strong>Built with using Go</strong>
</div>
