# User Management API

A robust RESTful API built with Go, featuring JWT authentication, role-based access control, and comprehensive user management capabilities.

##  Features

- **Authentication & Authorization**
  - User signup with email validation
  - Secure login with JWT tokens
  - Password strength validation
  - Bcrypt password hashing
  - HTTP-only secure cookies
  - Role-based access control (user/admin)

- **User Management**
  - CRUD operations for users
  - Dynamic age calculation from date of birth
  - Pagination support
  - Input validation

## Prerequisites

- Docker & Docker Compose
- Go 1.24+ (for local development)
- PostgreSQL 15 (handled by Docker)

## Tech Stack

- **Framework**: [Fiber](https://gofiber.io/) v2.52.10
- **Database**: PostgreSQL 15
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **Logging**: Zap
- **Validation**: go-playground/validator
- **SQL**: SQLC for type-safe queries
- **Testing**: Go testing package

## Installation

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/Sindoo24/User-API-Management.git
cd User-API_Management
```

2. Create environment file:
```bash
cp .env.example .env
```

3. Update `.env` with your configuration:
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable
SERVER_PORT=8080
LOG_LEVEL=info
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY_HOURS=24
COOKIE_SECURE=true
```

4. Start the application:
```bash
docker-compose up -d
```

5. Verify the application is running:
```bash
docker-compose logs -f api
```

