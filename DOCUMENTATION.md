# QStack-Backend - Complete Documentation

## Table of Contents
1. [Project Overview](#1-project-overview)
2. [Technology Stack](#2-technology-stack)
3. [Architecture](#3-architecture)
4. [Business Rules](#4-business-rules)
5. [API Documentation](#5-api-documentation)
6. [Database Schema](#6-database-schema)
7. [Configuration](#7-configuration)
8. [Queue System](#8-queue-system)

---

## 1. Project Overview

**QStack-Backend** is a Q&A platform backend (similar to Stack Overflow) built with Go. It provides a complete system for users to ask questions, provide answers, vote on content, and engage in technical discussions.

### Core Features
- User registration and authentication with JWT
- Email verification system
- Password reset functionality
- Question management with tags
- Answer system with acceptance marking
- Comment system on answers
- Voting system (upvote/downvote)
- Personalized question feeds
- User profiles with activity tracking
- Community statistics
- Async email processing via RabbitMQ

---

## 2. Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.25.5 |
| Web Framework | Echo | v4 |
| Database | PostgreSQL | - |
| ORM | GORM | v1.31.1 |
| Migrations | golang-migrate/migrate | v4.19.1 |
| Message Queue | RabbitMQ | amqp091-go |
| Email Testing | Mailpit | - |
| Authentication | JWT | golang-jwt/jwt/v5 |
| Validation | go-playground/validator | v10 |
| Password Hashing | bcrypt | golang.org/x/crypto |

---

## 3. Architecture

### Directory Structure
```
QStack-Backend/
├── cmd/                          # Application entry points
│   ├── server/main.go            # HTTP API server
│   ├── worker/main.go            # Background job worker
│   └── migrator/main.go          # Database migration CLI
├── internal/
│   ├── api/
│   │   ├── handlers/             # HTTP request handlers
│   │   ├── routes/               # Route definitions
│   │   └── middleware/           # JWT authentication
│   ├── config/                   # Configuration & DB connection
│   ├── models/
│   │   ├── domains/              # Database entities
│   │   └── dtos/                 # Request/Response DTOs
│   ├── repositories/             # Data access layer (GORM)
│   ├── services/                 # Business logic layer
│   ├── queue/                    # RabbitMQ producer/consumer
│   ├── workers/                  # Background job processors
│   └── validator/                # Request validation
├── migrations/                   # SQL migration files
├── docker/
│   └── docker-compose.yml        # RabbitMQ + Mailpit
├── uploads/                      # User uploaded files
└── go.mod, go.sum                # Dependencies
```

### Layer Architecture
```
┌─────────────────┐
│   HTTP Layer    │  (Echo Handlers)
├─────────────────┤
│  Business Layer │  (Services)
├─────────────────┤
│   Data Layer    │  (Repositories - GORM)
├─────────────────┤
│    Database     │  (PostgreSQL)
└─────────────────┘
       ↕
┌─────────────────┐
│  Message Queue  │  (RabbitMQ)
├─────────────────┤
│    Workers      │  (Email Worker)
└─────────────────┘
```

---

## 4. Business Rules

### 4.1 User Management

#### Registration (Signup)
- **Email must be unique** - Cannot register with an existing email
- **Username must be unique** - Cannot register with an existing username
- **Username minimum length**: 3 characters
- **Password minimum length**: 6 characters
- **Email format**: Must be valid email format
- **Email verification required**: New users must verify email before login
- **Password hashing**: All passwords are hashed using bcrypt with default cost
- **Email verification token**: Generated on signup, expires in 24 hours
- **Verification URL**: Returned in signup response for email verification

#### Login
- **Identifier**: Can be either email OR username
- **Password validation**: Compared against bcrypt hash
- **Email verification required**: Cannot login without verified email
- **Access token**: JWT, expires in 15 minutes
- **Refresh token**: JWT, expires in 7 days
- **Token storage**: HTTP-only cookies (XSS protection)
- **Invalid credentials**: Generic error message (security)

#### Email Verification
- **Token format**: 32-byte random, base64 URL-encoded
- **Token storage**: SHA256 hash stored in database
- **Token expiration**: 24 hours from creation
- **Single use**: Token marked as used after verification
- **Verification endpoint**: GET with query parameter `token`

#### Password Management
- **Change Password** (authenticated):
  - Must provide current password
  - New password minimum 6 characters
  - New password must be different from current
  - Requires JWT authentication

- **Forgot Password**:
  - Accepts email address
  - Does not reveal if email exists (security)
  - Generates reset token (24-hour expiry)
  - Sends reset email via RabbitMQ queue

- **Reset Password**:
  - Requires valid token
  - Token must not be expired
  - Token must not be used
  - New password minimum 6 characters
  - Token hashed with SHA256 for storage

#### Profile Management
- **Bio**: Optional text field, can be updated
- **Profile Image**: PNG format, stored in `/uploads/profile-images/`
- **Email notifications**: Configurable (default: disabled)
- **Activity tracking**: Questions, answers, votes, edits, accepted answers

### 4.2 Question Management

#### Creating Questions
- **Title**: Required, 10-200 characters
- **Description**: Required, minimum 20 characters
- **Tags**: Required, minimum 1 tag, each tag required
- **Tag processing**: 
  - Tags auto-created if don't exist
  - Tags normalized (lowercase, trimmed)
  - Maximum 50 characters per tag name
- **Author tracking**: Question linked to creating user
- **Initial vote count**: 0
- **Initial answer count**: 0

#### Updating Questions
- **Authorization**: Only question owner can update
- **Partial updates**: Any field can be updated independently
- **Timestamp**: UpdatedAt timestamp refreshed on change

#### Deleting Questions
- **Authorization**: Only question owner can delete
- **Cascade**: Related votes and tags are cascade-deleted

#### Question Feed
- **Public feed**: Available without authentication
- **Pagination**: Default 20 items, configurable limit/offset
- **Search**: By title (case-insensitive, partial match)
- **Tag filter**: Filter by specific tag
- **Sorting options**:
  - `date` (default): Newest first
  - `votes`: Highest voted first

#### Personalized Feed (My Feed)
- **Authentication required**: JWT token needed
- **Interest calculation**: Based on:
  - Tags from user's own questions
  - Tags from questions user voted on
  - Tags from questions user answered
- **Excludes**: User's own questions from feed
- **Pagination**: Default 20 items

#### User's Questions
- **Authentication required**: Can only view own questions
- **Pagination**: Default 20 items
- **Includes**: Full question data with tags and votes

### 4.3 Voting System

#### Voting Rules
- **Vote values**: +1 (upvote) or -1 (downvote) only
- **Owner restriction**: Cannot vote on own question
- **Authentication required**: Must be logged in
- **Vote uniqueness**: One vote per user per question

#### Vote Behavior
- **First vote**: Creates vote record, updates question vote count
- **Same vote (toggle)**: Removes vote, adjusts count accordingly
- **Changed vote**: Updates vote value, adjusts count by difference
- **Vote removal**: Deleting vote subtracts its value from count

### 4.4 Answer Management

#### Creating Answers
- **Description**: Required, minimum 10 characters
- **Question must exist**: Valid question ID required
- **Authentication required**: Must be logged in
- **Answer count**: Increments question's answer count

#### Updating Answers
- **Authorization**: Only answer owner can update
- **Partial updates**: Description can be updated
- **Timestamp**: UpdatedAt refreshed on change

#### Deleting Answers
- **Authorization**: Only answer owner can delete
- **Answer count**: Decrements question's answer count
- **Accepted status**: If deleted, question loses accepted answer

#### Accepting Answers
- **Authorization**: Only question owner can accept
- **Single accepted**: Database constraint ensures one accepted per question
- **Accepted display**: Accepted answers shown first

#### Answer Retrieval
- **By question**: All answers for a question
- **Ordering**: Accepted first, then by creation date (ascending)
- **Public access**: No authentication required

### 4.5 Comment System

#### Creating Comments
- **Parent type**: Only answers supported (type=2)
- **Body**: Required, 2-1000 characters
- **Authentication required**: Must be logged in
- **Answer must exist**: Valid answer ID required

#### Updating Comments
- **Authorization**: Only comment owner can update
- **Body**: Can be updated
- **No timestamp update**: CreatedAt remains original

#### Deleting Comments
- **Authorization**: Only comment owner can delete
- **Cascade**: No cascade (comments have no children)

#### Comment Retrieval
- **By answer**: All comments on an answer
- **Ordering**: Chronological (oldest first)
- **Public access**: No authentication required

### 4.6 Tags System

#### Tag Management
- **Auto-creation**: Tags created automatically with questions
- **Normalization**: Lowercase, trimmed whitespace
- **Uniqueness**: Tag names are unique
- **Maximum length**: 50 characters

#### Popular Tags
- **Limit**: Top 10 tags by default
- **Counting**: Based on question-tag associations
- **Public access**: No authentication required

### 4.7 User Activity & Statistics

#### Activity Tracking
- **Questions**: User's questions (last 5)
- **Answers**: User's answers (last 5)
- **Votes**: User's votes on questions (last 5)
- **Accepted answers**: User's accepted answers (last 5)
- **Edited questions**: Questions user modified (last 5)
- **Edited answers**: Answers user modified (last 5)
- **Sorting**: Most recent first
- **Privacy**: Users can only view their own activity

#### Profile Statistics
- **Total questions**: Count of user's questions
- **Total answers**: Count of user's answers
- **Total votes**: Sum of vote counts on user's questions
- **Preferred tags**: User's interested tags (future feature)

#### Community Statistics
- **Total users**: Count of all users
- **Total questions**: Count of all questions
- **Total answers**: Count of all answers
- **Public access**: No authentication required

#### User Listing
- **Pagination**: 20 users per page
- **Ordering**: Newest users first
- **Public data**: Username, bio, stats, join date

### 4.8 File Upload

#### Profile Image Upload
- **Format**: PNG (saved as .png)
- **Location**: `/uploads/profile-images/user-{id}.png`
- **Authentication required**: Must be logged in
- **Single image**: Overwrites previous image
- **Directory creation**: Auto-creates uploads folder

#### Image Upload (General)
- **Endpoint**: POST `/api/v1/upload`
- **No authentication**: Public endpoint
- **Form field**: `image`

### 4.9 Security Rules

#### Authentication
- **JWT tokens**: Signed with HS256
- **Secret**: Loaded from environment variable
- **Header format**: `Authorization: Bearer <token>`
- **Cookie fallback**: `access_token` cookie
- **Token extraction**: Header preferred over cookie

#### Authorization
- **Owner checks**: Users can only modify own content
- **Question owner privileges**: Can accept answers
- **Activity privacy**: Users can only view own activity
- **Vote ownership**: Cannot vote on own content

#### Data Protection
- **Password hashing**: bcrypt with default cost
- **Token hashing**: SHA256 for storage
- **HTTP-only cookies**: XSS protection
- **Generic errors**: Don't reveal user existence
- **Input validation**: All requests validated

#### Database Constraints
- **Unique email**: Enforced at database level
- **Unique username**: Enforced at database level
- **Unique votes**: One vote per user per question
- **Single accepted answer**: One per question
- **Cascade deletes**: Maintains referential integrity

### 4.10 Email System

#### Email Queue
- **Queue name**: `email_verification_queue`
- **Durable**: Survives RabbitMQ restarts
- **Manual ACK**: Jobs acknowledged after processing
- **Retry**: Failed jobs requeued

#### Email Types
- **Verification email**: Sent on signup
- **Password reset email**: Sent on forgot password request

#### Email Configuration
- **SMTP server**: Mailpit (testing)
- **No authentication**: Mailpit default
- **From address**: `no-reply@qstack.com`
- **Content type**: Plain text with UTF-8

---

## 5. API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

All protected endpoints require JWT authentication via:
- **Authorization Header**: `Authorization: Bearer <access_token>`
- **Cookie**: `access_token=<access_token>`

---

### 5.1 Authentication Endpoints

#### POST /auth/signup
Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "secret123"
}
```

**Validation Rules:**
- `email`: Required, valid email format
- `username`: Required, minimum 3 characters
- `password`: Required, minimum 6 characters

**Response (201 Created):**
```json
{
  "message": "Signup successful. Please verify your email.",
  "verify_url": "http://localhost:8080/verify-email?token=..."
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input or email already registered

---

#### POST /auth/login
Authenticate user and receive tokens.

**Request Body:**
```json
{
  "identifier": "johndoe@example.com",
  "password": "secret123"
}
```

**Validation Rules:**
- `identifier`: Required (email or username)
- `password`: Required

**Response (200 OK):**
```json
{
  "message": "login successful"
}
```

**Cookies Set:**
- `access_token`: JWT, expires in 1 day, HTTP-only
- `refresh_token`: JWT, expires in 7 days, HTTP-only

**Error Responses:**
- `401 Unauthorized`: Invalid credentials or email not verified

---

#### GET /auth/verify-email
Verify user's email address.

**Query Parameters:**
- `token`: Required, verification token from signup

**Response (200 OK):**
```json
{
  "message": "Email verified successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid or expired token

---

#### POST /auth/forgot-password
Request password reset email.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Validation Rules:**
- `email`: Required, valid email format

**Response (200 OK):**
```json
{
  "message": "If the email exists, reset link was sent"
}
```

**Note:** Always returns success to prevent email enumeration.

---

#### POST /auth/reset-password
Reset password with token.

**Request Body:**
```json
{
  "token": "reset-token-from-email",
  "new_password": "newsecret123"
}
```

**Validation Rules:**
- `token`: Required
- `new_password`: Required, minimum 6 characters

**Response (200 OK):**
```json
{
  "message": "password reset successful"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid or expired token

---

#### POST /auth/change-password
Change password (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Request Body:**
```json
{
  "current_password": "oldsecret",
  "new_password": "newsecret123"
}
```

**Validation Rules:**
- `current_password`: Required
- `new_password`: Required, minimum 6 characters

**Response (200 OK):**
```json
{
  "message": "password changed successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Incorrect current password

---

#### POST /auth/logout
Clear authentication cookies.

**Response (200 OK):**
```json
{
  "message": "logged out"
}
```

---

### 5.2 Question Endpoints

#### GET /questions
Get public question feed.

**Query Parameters:**
- `search`: Optional, search in title
- `tag`: Optional, filter by tag name
- `sort`: Optional, `date` (default) or `votes`
- `limit`: Optional, default 20
- `offset`: Optional, default 0

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "title": "How to use Go generics?",
    "description": "I'm trying to understand...",
    "vote_count": 5,
    "answer_count": 2,
    "author": {
      "id": 1,
      "username": "gopher"
    },
    "tags": ["go", "generics"],
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

---

#### GET /questions/:id
Get single question by ID.

**Path Parameters:**
- `id`: Question ID

**Response (200 OK):**
```json
{
  "id": 1,
  "title": "How to use Go generics?",
  "description": "I'm trying to understand...",
  "vote_count": 5,
  "answer_count": 2,
  "author": {
    "id": 1,
    "username": "gopher"
  },
  "tags": ["go", "generics"],
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

**Error Responses:**
- `404 Not Found`: Question not found

---

#### GET /questions/my-feed
Get personalized question feed (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Query Parameters:**
- `limit`: Optional, default 20
- `offset`: Optional, default 0

**Response (200 OK):**
```json
[
  {
    "id": 2,
    "title": "Understanding channels in Go",
    "description": "Can someone explain...",
    "vote_count": 3,
    "answer_count": 1,
    "author": {
      "id": 2,
      "username": "concurrency_fan"
    },
    "tags": ["go", "channels"],
    "created_at": "2024-01-02T10:00:00Z",
    "updated_at": "2024-01-02T10:00:00Z"
  }
]
```

**Note:** Returns questions based on user's tag interests, excludes own questions.

---

#### GET /questions/my
Get current user's questions (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Query Parameters:**
- `limit`: Optional, default 20
- `offset`: Optional, default 0

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "title": "How to use Go generics?",
    "description": "I'm trying to understand...",
    "vote_count": 5,
    "answer_count": 2,
    "author": {
      "id": 1,
      "username": "johndoe"
    },
    "tags": ["go", "generics"],
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

---

#### POST /questions
Create new question (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Request Body:**
```json
{
  "title": "How to use Go generics?",
  "description": "I'm trying to understand how generics work in Go 1.18+...",
  "tags": ["go", "generics", "types"]
}
```

**Validation Rules:**
- `title`: Required, 10-200 characters
- `description`: Required, minimum 20 characters
- `tags`: Required, minimum 1 tag, each tag required

**Response (201 Created):**
```json
{
  "id": 1,
  "title": "How to use Go generics?",
  "description": "I'm trying to understand how generics work in Go 1.18+...",
  "vote_count": 0,
  "answer_count": 0,
  "author": {
    "id": 1,
    "username": "johndoe"
  },
  "tags": ["go", "generics", "types"],
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

---

#### PUT /questions/:id
Update question (authenticated, owner only).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `id`: Question ID

**Request Body:**
```json
{
  "title": "Updated title",
  "description": "Updated description",
  "tags": ["go", "updated"]
}
```

**Note:** All fields optional, partial update supported.

**Response (200 OK):**
```json
{
  "message": "updated"
}
```

**Error Responses:**
- `400 Bad Request`: Not authorized or invalid input

---

#### DELETE /questions/:id
Delete question (authenticated, owner only).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `id`: Question ID

**Response (200 OK):**
```json
{
  "message": "deleted"
}
```

**Error Responses:**
- `403 Forbidden`: Not authorized

---

#### POST /questions/:id/vote
Vote on question (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `id`: Question ID

**Request Body:**
```json
{
  "value": 1
}
```

**Validation Rules:**
- `value`: Required, must be 1 or -1

**Response (200 OK):**
```json
{
  "message": "vote updated"
}
```

**Error Responses:**
- `400 Bad Request`: Cannot vote on own question or invalid value

---

### 5.3 Answer Endpoints

#### GET /answers/question/:question_id
Get all answers for a question.

**Path Parameters:**
- `question_id`: Question ID

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "description": "Here's how to use generics...",
    "is_accepted": true,
    "author": {
      "id": 2,
      "username": "expert"
    },
    "created_at": "2024-01-01T11:00:00Z",
    "updated_at": "2024-01-01T11:00:00Z"
  }
]
```

**Note:** Accepted answers returned first, then by creation date.

---

#### POST /answers/question/:question_id
Create answer (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `question_id`: Question ID

**Request Body:**
```json
{
  "description": "Here's how you can solve this problem..."
}
```

**Validation Rules:**
- `description`: Required, minimum 10 characters

**Response (201 Created):**
```json
{
  "id": 1,
  "description": "Here's how you can solve this problem...",
  "is_accepted": false,
  "author": {
    "id": 1,
    "username": "johndoe"
  },
  "created_at": "2024-01-01T11:00:00Z",
  "updated_at": "2024-01-01T11:00:00Z"
}
```

---

#### PUT /answers/:id
Update answer (authenticated, owner only).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `id`: Answer ID

**Request Body:**
```json
{
  "description": "Updated answer description"
}
```

**Response (200 OK):**
```json
{
  "message": "updated"
}
```

**Error Responses:**
- `403 Forbidden`: Not authorized

---

#### DELETE /answers/:id
Delete answer (authenticated, owner only).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `id`: Answer ID

**Response (200 OK):**
```json
{
  "message": "deleted"
}
```

**Error Responses:**
- `403 Forbidden`: Not authorized

---

#### PUT /answers/:id/accept
Accept answer (authenticated, question owner only).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `id`: Answer ID

**Response (200 OK):**
```json
{
  "message": "accepted"
}
```

**Error Responses:**
- `403 Forbidden`: Only question owner can accept answer

---

### 5.4 Comment Endpoints

#### GET /comments/answer/:answer_id
Get all comments on an answer.

**Path Parameters:**
- `answer_id`: Answer ID

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "body": "Great answer! This helped me a lot.",
    "author": {
      "id": 3,
      "username": "thankful_user"
    },
    "created_at": "2024-01-01T12:00:00Z"
  }
]
```

---

#### POST /comments/answer/:answer_id
Create comment (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `answer_id`: Answer ID

**Request Body:**
```json
{
  "body": "Thanks for this detailed explanation!"
}
```

**Validation Rules:**
- `body`: Required, 2-1000 characters

**Response (201 Created):**
```json
{
  "id": 1,
  "body": "Thanks for this detailed explanation!",
  "author": {
    "id": 1,
    "username": "johndoe"
  },
  "created_at": "2024-01-01T12:00:00Z"
}
```

---

#### PUT /comments/:id
Update comment (authenticated, owner only).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `id`: Comment ID

**Request Body:**
```json
{
  "body": "Updated comment text"
}
```

**Response (200 OK):**
```json
{
  "message": "updated"
}
```

**Error Responses:**
- `403 Forbidden`: Not authorized

---

#### DELETE /comments/:id
Delete comment (authenticated, owner only).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `id`: Comment ID

**Response (200 OK):**
```json
{
  "message": "deleted"
}
```

**Error Responses:**
- `403 Forbidden`: Not authorized

---

### 5.5 User Endpoints

#### GET /users
Get all users (paginated).

**Query Parameters:**
- `page`: Optional, default 1

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "username": "johndoe",
    "bio": "Go enthusiast",
    "total_questions": 5,
    "total_answers": 10,
    "total_votes": 25,
    "created_at": "2024-01-01T10:00:00Z"
  }
]
```

---

#### GET /users/:id/profile
Get user profile by ID.

**Path Parameters:**
- `id`: User ID

**Response (200 OK):**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "johndoe@example.com",
  "bio": "Go enthusiast",
  "profile_image": "/uploads/profile-images/user-1.png",
  "total_questions": 5,
  "total_answers": 10,
  "total_votes": 25,
  "preferred_tags": ["go", "backend"],
  "created_at": "2024-01-01T10:00:00Z"
}
```

**Error Responses:**
- `404 Not Found`: User not found

---

#### GET /users/me
Get current user profile (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Response (200 OK):**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "johndoe@example.com",
  "bio": "Go enthusiast",
  "profile_image": "/uploads/profile-images/user-1.png",
  "total_questions": 5,
  "total_answers": 10,
  "total_votes": 25,
  "preferred_tags": ["go", "backend"],
  "created_at": "2024-01-01T10:00:00Z"
}
```

---

#### PUT /users/profile
Update user profile (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Request Body:**
```json
{
  "bio": "Updated bio text"
}
```

**Response (200 OK):**
```json
{
  "message": "profile updated"
}
```

---

#### POST /users/profile/image
Upload profile image (authenticated).

**Headers:**
- `Authorization: Bearer <access_token>`

**Request:**
- Content-Type: `multipart/form-data`
- Form field: `image` (file)

**Response (200 OK):**
```json
{
  "profile_image": "/uploads/profile-images/user-1.png"
}
```

**Error Responses:**
- `400 Bad Request`: Image required

---

#### GET /users/:id/activity
Get user activity (authenticated, own activity only).

**Headers:**
- `Authorization: Bearer <access_token>`

**Path Parameters:**
- `id`: User ID (must match authenticated user)

**Response (200 OK):**
```json
[
  {
    "type": "question",
    "title": "How to use Go generics?",
    "target_id": 1,
    "created_at": "2024-01-01T10:00:00Z"
  },
  {
    "type": "answer",
    "target_id": 5,
    "created_at": "2024-01-01T11:00:00Z"
  },
  {
    "type": "vote",
    "target_id": 3,
    "value": 1,
    "created_at": "2024-01-01T12:00:00Z"
  },
  {
    "type": "accept",
    "target_id": 2,
    "created_at": "2024-01-01T13:00:00Z"
  },
  {
    "type": "edit",
    "target_id": 1,
    "created_at": "2024-01-01T14:00:00Z"
  }
]
```

**Activity Types:**
- `question`: User asked a question
- `answer`: User answered a question
- `vote`: User voted on a question
- `accept`: User's answer was accepted
- `edit`: User edited a question or answer

**Error Responses:**
- `403 Forbidden`: Can only view own activity

---

#### GET /users/community/stats
Get community statistics.

**Response (200 OK):**
```json
{
  "total_users": 150,
  "total_questions": 500,
  "total_answers": 1200
}
```

---

### 5.6 Tag Endpoints

#### GET /tags/popular
Get popular tags.

**Response (200 OK):**
```json
[
  {
    "tag": "go",
    "count": 45
  },
  {
    "tag": "javascript",
    "count": 38
  },
  {
    "tag": "python",
    "count": 32
  }
]
```

---

### 5.7 Upload Endpoints

#### POST /upload
Upload image file.

**Request:**
- Content-Type: `multipart/form-data`
- Form field: `image` (file)

**Response (200 OK):**
```json
{
  "url": "/uploads/image-123.png"
}
```

---

### 5.8 Health Check

#### GET /health
Check API and database health.

**Response (200 OK):**
```json
{
  "status": "ok",
  "db": "connected"
}
```

**Error Responses:**
- `500 Internal Server Error`: Database connection failed

---

### 5.9 Protected Test Route

#### GET /api/v1/protected/me
Test JWT authentication.

**Headers:**
- `Authorization: Bearer <access_token>`

**Response (200 OK):**
```json
{
  "message": "This is a protected route",
  "user_id": 1
}
```

---

## 6. Database Schema

### Tables

#### users
| Column | Type | Constraints |
|--------|------|-------------|
| id | BIGSERIAL | PRIMARY KEY, AUTO INCREMENT |
| email | VARCHAR(255) | NOT NULL, UNIQUE |
| password_hash | TEXT | NOT NULL |
| username | VARCHAR(50) | NOT NULL, UNIQUE |
| bio | TEXT | NULL |
| profile_image | TEXT | NULL |
| email_notifications_enabled | BOOLEAN | NOT NULL, DEFAULT FALSE |
| email_verified | BOOLEAN | NOT NULL, DEFAULT FALSE |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

#### tags
| Column | Type | Constraints |
|--------|------|-------------|
| id | BIGSERIAL | PRIMARY KEY, AUTO INCREMENT |
| name | VARCHAR(50) | NOT NULL, UNIQUE |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

#### user_preferred_tags
| Column | Type | Constraints |
|--------|------|-------------|
| user_id | BIGINT | NOT NULL, FK → users(id), PRIMARY KEY |
| tag_id | BIGINT | NOT NULL, FK → tags(id), PRIMARY KEY |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

#### user_tag_interest
| Column | Type | Constraints |
|--------|------|-------------|
| user_id | BIGINT | NOT NULL, FK → users(id), PRIMARY KEY |
| tag_id | BIGINT | NOT NULL, FK → tags(id), PRIMARY KEY |
| score | INT | NOT NULL, DEFAULT 0 |
| updated_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

#### questions
| Column | Type | Constraints |
|--------|------|-------------|
| id | BIGSERIAL | PRIMARY KEY, AUTO INCREMENT |
| user_id | BIGINT | NOT NULL, FK → users(id) |
| title | VARCHAR(200) | NOT NULL |
| description | TEXT | NOT NULL |
| vote_count | INT | NOT NULL, DEFAULT 0 |
| answer_count | INT | NOT NULL, DEFAULT 0 |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Indexes:**
- `idx_questions_user_id` ON user_id
- `idx_questions_created_at` ON created_at DESC
- `idx_questions_vote_count` ON vote_count DESC
- `idx_questions_title` ON title

#### question_tags
| Column | Type | Constraints |
|--------|------|-------------|
| question_id | BIGINT | NOT NULL, FK → questions(id), PRIMARY KEY |
| tag_id | BIGINT | NOT NULL, FK → tags(id), PRIMARY KEY |

**Indexes:**
- `idx_question_tags_tag_id_question_id` ON (tag_id, question_id)

#### answers
| Column | Type | Constraints |
|--------|------|-------------|
| id | BIGSERIAL | PRIMARY KEY, AUTO INCREMENT |
| question_id | BIGINT | NOT NULL, FK → questions(id) |
| user_id | BIGINT | NOT NULL, FK → users(id) |
| description | TEXT | NOT NULL |
| is_accepted | BOOLEAN | NOT NULL, DEFAULT FALSE |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Indexes:**
- `idx_answers_question_id_created_at` ON (question_id, created_at)
- `ux_answers_one_accepted_per_question` UNIQUE ON (question_id) WHERE is_accepted = TRUE

#### comments
| Column | Type | Constraints |
|--------|------|-------------|
| id | BIGSERIAL | PRIMARY KEY, AUTO INCREMENT |
| user_id | BIGINT | NOT NULL, FK → users(id) |
| parent_type | SMALLINT | NOT NULL (1=question, 2=answer, 3=comment) |
| parent_id | BIGINT | NOT NULL |
| body | VARCHAR(1000) | NOT NULL |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Indexes:**
- `idx_comments_parent` ON (parent_type, parent_id)

#### question_votes
| Column | Type | Constraints |
|--------|------|-------------|
| id | BIGSERIAL | PRIMARY KEY, AUTO INCREMENT |
| question_id | BIGINT | NOT NULL, FK → questions(id) |
| user_id | BIGINT | NOT NULL, FK → users(id) |
| value | SMALLINT | NOT NULL, CHECK (value IN (1, -1)) |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Constraints:**
- `ux_question_votes_unique` UNIQUE (question_id, user_id)

**Indexes:**
- `idx_question_votes_question_id` ON question_id
- `idx_question_votes_user_id` ON user_id

#### email_verification_tokens
| Column | Type | Constraints |
|--------|------|-------------|
| id | BIGSERIAL | PRIMARY KEY, AUTO INCREMENT |
| user_id | BIGINT | NOT NULL, FK → users(id) |
| token_hash | TEXT | NOT NULL |
| expires_at | TIMESTAMPTZ | NOT NULL |
| used_at | TIMESTAMPTZ | NULL |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Indexes:**
- `idx_email_verification_token_hash` ON token_hash

#### password_reset_tokens
| Column | Type | Constraints |
|--------|------|-------------|
| id | BIGSERIAL | PRIMARY KEY, AUTO INCREMENT |
| user_id | BIGINT | NOT NULL, FK → users(id) |
| token_hash | TEXT | NOT NULL |
| expires_at | TIMESTAMPTZ | NOT NULL |
| used_at | TIMESTAMPTZ | NULL |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Indexes:**
- `idx_password_reset_token_hash` ON token_hash

#### notifications
| Column | Type | Constraints |
|--------|------|-------------|
| id | BIGSERIAL | PRIMARY KEY, AUTO INCREMENT |
| user_id | BIGINT | NOT NULL, FK → users(id) |
| actor_user_id | BIGINT | NULL, FK → users(id) |
| type | SMALLINT | NOT NULL (1=answer, 2=comment, 3=reply) |
| entity_type | SMALLINT | NOT NULL (1=question, 2=answer, 3=comment) |
| entity_id | BIGINT | NOT NULL |
| is_read | BOOLEAN | NOT NULL, DEFAULT FALSE |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |
| sent_email_at | TIMESTAMPTZ | NULL |

**Indexes:**
- `idx_notifications_user_unread_created` ON (user_id, is_read, created_at DESC)

---

## 7. Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | `8080` | HTTP server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD` | *(required)* | Database password |
| `DB_NAME` | `qstack` | Database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `JWT_SECRET` | *(required)* | JWT signing secret |
| `APP_BASE_URL` | `http://localhost:8080` | Base URL for email links |
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:5672/` | RabbitMQ connection URL |
| `MAILPIT_HOST` | `localhost` | Mailpit SMTP host |
| `MAILPIT_PORT` | `1025` | Mailpit SMTP port |

### Database Connection
- **Timezone**: Asia/Dhaka
- **Max open connections**: 25
- **Max idle connections**: 10
- **Connection lifetime**: 5 minutes

---

## 8. Queue System

### RabbitMQ Configuration

#### Queue: `email_verification_queue`
- **Durable**: Yes (survives restarts)
- **Auto-delete**: No
- **Exclusive**: No
- **Acknowledgment**: Manual

### Job Structure

```json
{
  "email": "user@example.com",
  "token": "verification-or-reset-token",
  "type": "verify or reset"
}
```

### Job Types

#### Verification Email (`type: "verify"`)
- **Trigger**: User signup
- **Subject**: "Verify Your Email"
- **Link**: `{APP_BASE_URL}/verify-email?token={token}`

#### Password Reset Email (`type: "reset"`)
- **Trigger**: Forgot password request
- **Subject**: "Reset Your Password"
- **Link**: `{APP_BASE_URL}/reset-password?token={token}`

### Worker Process

1. **Consume**: Worker listens on `email_verification_queue`
2. **Deserialize**: Parse JSON job payload
3. **Process**: Send email via SMTP to Mailpit
4. **Acknowledge**: Manual ACK on success, NACK on failure
5. **Retry**: Failed jobs requeued for retry

### Email Configuration
- **From**: `no-reply@qstack.com`
- **SMTP Host**: Mailpit (no authentication)
- **Content-Type**: `text/plain; charset=UTF-8`

---

## Appendix: Error Codes

### HTTP Status Codes

| Code | Meaning | Common Scenarios |
|------|---------|------------------|
| 200 | OK | Successful request |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid input, validation errors |
| 401 | Unauthorized | Missing or invalid authentication |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server error, database issues |

### Common Error Messages

- `"invalid credentials"` - Wrong email/username or password
- `"email not verified"` - Login attempted before email verification
- `"not authorized"` - Trying to modify someone else's content
- `"cannot vote own question"` - Question owner trying to vote
- `"only question owner can accept answer"` - Wrong user accepting answer
- `"invalid or expired token"` - Token verification failed
- `"email already registered"` - Duplicate email on signup
- `"username already taken"` - Duplicate username on signup

---

## Appendix: Data Validation Summary

### User Validation
| Field | Rules |
|-------|-------|
| Email | Required, valid email format, unique |
| Username | Required, min 3 chars, unique |
| Password | Required, min 6 chars |
| Bio | Optional, text |

### Question Validation
| Field | Rules |
|-------|-------|
| Title | Required, 10-200 chars |
| Description | Required, min 20 chars |
| Tags | Required, min 1 tag, each required |

### Answer Validation
| Field | Rules |
|-------|-------|
| Description | Required, min 10 chars |

### Comment Validation
| Field | Rules |
|-------|-------|
| Body | Required, 2-1000 chars |

### Vote Validation
| Field | Rules |
|-------|-------|
| Value | Required, must be 1 or -1 |

---

*Documentation generated from QStack-Backend codebase*
