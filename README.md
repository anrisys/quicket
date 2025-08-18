# Quicket API

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An event and booking management API built in Go, designed to demonstrate a modular monolith architecture with a clear path to microservices. This project includes a payments simulation, authentication, and role-based access control.

## 🚀 Features

- User Management: Public endpoints for user registration and login.

- Event Management: Authenticated and role-protected endpoints for creating events.

- Booking System: Authenticated endpoints for booking events.

- Payments Simulation: A dedicated payment simulation endpoint to demonstrate a full transaction flow.

- Authentication & Authorization: Utilizes JWT for secure access to protected routes with role-based checks.

## 🛠️ Technologies Used

- Go: The primary language for the backend.

- Gin: A high-performance HTTP web framework.

- GORM: The ORM library for database interactions.

- Viper: A library for handling configuration.

- Wire: A dependency injection tool to manage application components.

- JWT Auth: For authentication

- Bcrypt: For secure password hashing.

- ZeroLog: A lightweight and fast logging library.

- Golang-migrate: For database schema migrations.

- MySQL: The primary database.

- Swagger : For API docs

## 📁 Project Structure

The project is structured following the principles of a modular monolith, which allows for a clean separation of business logic and a clear path for future migration to a microservices architecture.

- cmd/server/: The application's entry point.

- internal/: Contains all private, application-specific code. This includes business logic for booking, event, payment, and user.

- pkg/: Houses shared, reusable packages and utilities that can be imported by other services or projects.

- migration/: Stores database migration scripts.

- api/docs/: Location for generated API documentation.

  ```
  quicket/
  ├──  api/               # API contracts or Swagger/OpenAPI
  │ └── docs/             # Generated API docs (Swagger)
  ├── cmd/                # Application entrypoints
  │ └── server/
  │ └── main.go           # Main app entrypoint
  ├── internal/           # Business logic (domain-driven design)
  │ ├── booking/          # Booking domain (handler, service, repo)
  │ ├── dto/              # Request/response DTOs
  │ ├── event/            # Event domain
  │ ├── payment/          # Payment domain
  │ ├── user/             # User domain
  │ └── validations/      # Custom input validation
  ├── migration/          # Database migration files
  ├── pkg/                # Shared libraries/utilities
  │ ├── config/           # Viper-based config loader
  │ ├── database/         # Database connection + GORM
  │ ├── di/               # Dependency injection with Wire
  │ ├── middleware/       # Gin middlewares (JWT, roles)
  │ ├── security/         # Password hashing
  │ ├── token/            # JWT utilities
  │ ├── types/            # Shared enums/types
  │ └── util/             # Helper utilities
  ├── .env                # Local environment variables
  ├── .example.env        # Example env file
  ├── go.mod
  └── go.sum
  ```

## ▶️ Getting Started

### Prerequisites

- Go v1.21 or newer

- MySQL Server

- [Golang Migrate CLI](https://github.com/golang-migrate/migrate)

### Local Setup

1. Clone the repository:

   ```
   git clone https://github.com/your-username/quicket.git
   cd quicket
   ```

2. Configure environment variables:

   Copy the .example.env file and rename it to .env. Fill in your MySQL database credentials and a secure JWT secret.

   ```
   cp .example.env .env
   ```

3. Run database migrations (using golang migrate cli):
   ```
   migrate -database YOUR_DATABASE_URL -path PATH_TO_YOUR_MIGRATIONS up
   ```
   - Make sure you already have the database. If you use mysql server, you can rewrite the database URL in the .example.env.
   - For migration path, you can use `migration` since it is in the root folder
4. Install dependencies and run the server:
   ```
   go mod tidy
   go run ./cmd/server
   ```
   The server should now be running on http://localhost:8080.

## 📄 API Documentation

This project uses swag to automatically generate API documentation. To view the documentation, you must first generate it and run the server.

Detailed API Documentation

## 👤 Endpoints

Public

- `POST /api/v1/register`: Register a new user.

- `POST /api/v1/login`: Log in and receive a JWT token.

Protected (Requires JWT)

- `POST /api/v1/events`: Create a new event. (Roles: admin, organizer)

- `POST /api/v1/bookings`: Create a new booking. (All authenticated users)

- `POST /api/v1/payments`: Simulate a payment. (All authenticated users)

## 🔮 Next Steps

- Split into microservices (user-service, event-service, booking-service, payment-service)

- Add Redis for caching / async jobs

- Use API Gateway (Kong)

## License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
