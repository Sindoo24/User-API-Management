# Simple Backend Project

A RESTful API for managing users with date of birth and dynamically calculated age.

## Setup and Run

### Prerequisites
- Docker
- Docker Compose

### Running the Application

1. Clone the repository and navigate to the project directory

2. Start the application:
```bash
docker compose up --build -d
```

3. The API will be available at `http://localhost:8080`

### Stopping the Application
```bash
docker compose down
```

## API Usage

### Create User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```

### Get User (with calculated age)
```bash
curl http://localhost:8080/users/1
```

Response:
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}
```

### List All Users
```bash
curl http://localhost:8080/users
```

### Update User
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/users/1
```

## Features

- User management with CRUD operations
- Age automatically calculated from date of birth
- Pagination support: `GET /users?page=1&limit=10`
- Input validation (name minimum 2 characters, valid date format)
