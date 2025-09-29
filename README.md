# Cortex API

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Release](https://img.shields.io/github/v/release/CodeAndCraft-Online/cortex-api?color=blue)
[![Test Coverage](https://codecov.io/gh/CodeAndCraft-Online/cortex-api/branch/main/graph/badge.svg)](https://codecov.io/gh/CodeAndCraft-Online/cortex-api)
[![CI](https://github.com/CodeAndCraft-Online/cortex-api/actions/workflows/test-coverage.yml/badge.svg)](https://github.com/CodeAndCraft-Online/cortex-api/actions/workflows/test-coverage.yml)

A Reddit-like social media platform backend API built with Go and PostgreSQL, providing a comprehensive platform for community-driven discussions with posts, comments, voting systems, and private/public community features.

## 🚀 Features

### Core Functionality
- **User Authentication** - JWT-based authentication with refresh token rotation
- **Community Management** - Public and private communities (subs) with invitation system
- **Content Creation** - Posts with image support and threaded comments
- **Voting System** - Upvote/downvote functionality for posts and comments
- **Password Reset** - Secure token-based password recovery system

### Technical Highlights
- **Clean Architecture** - Layered approach with clear separation of concerns
- **RESTful API** - Consistent HTTP endpoints with proper status codes
- **Database Design** - PostgreSQL with GORM ORM and strategic indexing
- **Security First** - bcrypt password hashing, JWT tokens, CORS, rate limiting
- **Docker Support** - Containerized deployment with multi-stage builds

## 🛠 Technology Stack

- **Language**: Go 1.x
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 13+ with GORM ORM
- **Authentication**: JWT with refresh tokens
- **Containerization**: Docker
- **Development**: Hot reload, environment configuration

## 📁 Project Structure

```
cortex-api/
├── main.go                     # Application entry point
├── create_database.sql         # PostgreSQL schema DDL
├── Dockerfile                  # Multi-stage container build
├── .dockerignore
├── .gitignore
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency lock file
│
├── internal/
│   ├── database/
│   │   └── database.go         # Database initialization
│   │
│   ├── models/                 # Data structures and ORM models
│   │   ├── user.go             # User entity
│   │   ├── post.go             # Post entity
│   │   ├── comment.go          # Comment entity
│   │   ├── sub.go              # Community (sub) entity
│   │   ├── vote.go             # Vote entity
│   │   └── reset_token.go      # Password reset tokens
│   │
│   ├── repositories/           # Data access layer
│   │   ├── auth_repository.go
│   │   ├── post_repository.go
│   │   ├── comment_repository.go
│   │   ├── sub_repository.go
│   │   └── user_repository.go  # Repository implementations
│   │
│   ├── services/               # Business logic layer
│   │   ├── auth_service.go
│   │   ├── post_service.go
│   │   ├── comments_service.go
│   │   ├── sub_service.go
│   │   └── user_service.go     # Service implementations
│   │
│   ├── handlers/               # HTTP request/response handlers
│   │   ├── auth_handlers.go
│   │   ├── post_handlers.go
│   │   ├── votes_handlers.go
│   │   ├── sub_handlers.go
│   │   ├── user-auth_handlers.go
│   │   └── user-login.go       # Handler implementations
│   │
│   └── routes/                 # API routing and middleware
│       ├── routes.go           # Main routing setup
│       └── auth/, comments/, posts/, subs/, users/, votes/
│           └── *.go            # Route group definitions
│
├── pkg/                        # Shared utilities and middleware
│   ├── auth.go                 # Authentication utilities
│   └── rate_limit.go           # Rate limiting middleware
│
├── memory-bank/               # Project documentation (see .clinerules/)
│   ├── projectbrief.md        # Core requirements and goals
│   ├── productContext.md      # Why product exists, UX goals
│   ├── systemPatterns.md      # Architecture and design patterns
│   ├── techContext.md         # Technology stack and setup
│   ├── activeContext.md       # Current development focus
│   └── progress.md            # Development progress tracking
│
└── .vscode/
    └── .env                   # Environment variables (development)
```

## 🏗 Architecture

The API follows Clean Architecture principles with clear separation of concerns:

```
┌─────────────────┐
│   Routes Layer  │  ← Gin Router Groups
├─────────────────┤
│  Handlers Layer │  ← HTTP Request/Response
├─────────────────┤
│ Services Layer  │  ← Business Logic & Validation
├─────────────────┤
│Repository Layer │  ← Database Operations
├─────────────────┤
│   Models Layer  │  ← Data Transfer Objects
└─────────────────┘
```

### Key Design Patterns
- **Repository Pattern**: Abstract database operations
- **Service Layer**: Business logic orchestration
- **DTO Pattern**: Separate internal/external representations
- **Dependency Injection**: Clean component coupling

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Docker (optional, for containerized deployment)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/CodeAndCraft-Online/cortex-api.git
   cd cortex-api
   ```

2. **Environment Configuration**
   ```bash
   cp .vscode/.env .env
   # Edit .env with your database credentials
   ```

3. **Database Setup**
   ```bash
   # Ensure PostgreSQL is running locally
   createdb cortex_db
   psql -d cortex_db -f create_database.sql
   ```

4. **Install Dependencies**
   ```bash
   go mod download
   ```

5. **Run the Application**
   ```bash
   go run main.go
   # API will be available at http://localhost:8080
   ```

### Docker Deployment

```bash
# Build and run with Docker Compose
docker build -t cortex-api .
docker run -p 8080:8080 --env-file .env cortex-api
```

## 📡 API Endpoints

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh access token
- `POST /auth/reset-password` - Request password reset

### Users
- `GET /users/:id` - Get user profile
- `PUT /users/:id` - Update user profile

### Communities (Subs)
- `GET /subs` - List public communities
- `GET /subs/:id` - Get community details
- `POST /subs` - Create new community
- `PUT /subs/:id` - Update community
- `POST /subs/:id/join` - Join community
- `POST /subs/:id/invite` - Invite user to private community

### Posts
- `GET /posts` - List posts (with pagination)
- `GET /posts/:id` - Get specific post
- `POST /posts` - Create new post
- `PUT /posts/:id` - Update post
- `DELETE /posts/:id` - Delete post

### Comments
- `GET /posts/:id/comments` - Get post comments
- `POST /posts/:id/comments` - Create comment
- `PUT /comments/:id` - Update comment
- `DELETE /comments/:id` - Delete comment

### Voting
- `POST /posts/:id/vote` - Vote on post
- `POST /comments/:id/vote` - Vote on comment
- `DELETE /votes/:id` - Remove vote

## 🔧 Development Status

### ✅ Completed Features
- Complete PostgreSQL database schema
- JWT authentication system
- User registration and login
- Basic community (sub) management
- Post creation and retrieval
- Comment system foundation
- Voting system foundation

### 🔄 In Progress
- Complete comment CRUD operations
- Full voting system implementation
- Private community invitation system
- Password reset flow completion

### 📋 Planned Features
- Image upload handling
- Rate limiting expansion
- API documentation (OpenAPI/Swagger)
- Comprehensive testing suite
- Email notifications
- Caching layer implementation

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Clean Architecture principles
- Write comprehensive tests
- Maintain consistent error handling
- Use proper Go naming conventions
- Document complex business logic

## 🧪 Testing

### Test Coverage
The project includes comprehensive unit and integration tests with the following coverage areas:

- **Repository Layer**: Database operations and data access
- **Service Layer**: Business logic and validation
- **Handler Layer**: HTTP request/response handling
- **Integration Tests**: Full API endpoint testing

### Dependencies
- [testify](https://github.com/stretchr/testify) - Assertions and test utilities
- [dockertest](https://github.com/ory/dockertest/v3) - Integration test database setup
- [sqlmock](https://github.com/DATA-DOG/go-sqlmock) - SQL mocking for unit tests

### Running Tests

#### Local Development (requires PostgreSQL)
```bash
# Using the test script
./scripts/run-tests.sh

# Or manually
go test -v ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

#### With Docker (for integration tests)
```bash
# Build and run tests in container
docker build -t cortex-api-test .
docker run --rm cortex-api-test go test ./...
```

### CI/CD Coverage
- GitHub Actions workflow runs on every push/PR
- PostgreSQL service for integration testing
- Codecov integration for coverage tracking
- Coverage reports are generated and stored as artifacts

### Test Structure
```
tests/
├── unit/               # Unit tests (mocked dependencies)
├── integration/        # Integration tests (real DB)
└── coverage/           # Test coverage reports
```

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📞 Support

For support, create an issue in the GitHub repository.

---

**Built with ❤️ using Go, Gin, and PostgreSQL**
