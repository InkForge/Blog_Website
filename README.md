# Blog_Website Backend (Go, Gin, MongoDB, Clean Architecture)

## Overview
This is a robust, production-ready backend for a blogging platform, built with Go, Gin, MongoDB, and Clean Architecture principles. It features modular domain-driven design, JWT authentication, role-based access control, advanced search/filtering, and full CRUD for blogs, users, comments, tags, and reactions.

---

## Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Setup & Installation](#setup--installation)
- [Configuration](#configuration)
- [Running the App](#running-the-app)
- [API Endpoints](#api-endpoints)
- [Authentication & Roles](#authentication--roles)
- [Testing](#testing)
- [Development Workflow](#development-workflow)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

---

## Features
- Clean Architecture: strict separation of domain, usecase, repository, and delivery layers
- MongoDB with repository pattern and transaction support
- JWT authentication with role-based middleware (USER, ADMIN)
- Blog CRUD, search, filter, view count, like/dislike, and comment support
- Tag normalization and auto-creation
- User registration, login, OAuth2, password reset, email verification
- Integration and unit tests (testify)
- Configurable via `.env` or `config.env` (Viper)

---

## Architecture
- **domain/**: Pure business logic, interfaces, and error definitions
- **usecases/**: Application logic, orchestrates domain and repositories
- **repositories/**: MongoDB implementations, DTO conversion, data access
- **delivery/controllers/**: Gin HTTP handlers, request/response DTOs
- **delivery/routes/**: Route registration, middleware wiring
- **infrastructures/**: Auth, config, and external service integrations

---

## Tech Stack
- Go 1.20+
- Gin (HTTP framework)
- MongoDB (with replica set for transactions)
- Viper (config)
- Testify (testing)
- JWT (github.com/golang-jwt/jwt)

---

## Setup & Installation

### Prerequisites
- Go 1.20 or newer
- MongoDB (run as a replica set for transaction support)
- [Optional] MongoDB Compass for GUI

### Clone the Repository
```sh
git clone https://github.com/InkForge/Blog_Website.git
cd Blog_Website
```

### Install Dependencies
```sh
go mod tidy
```

### MongoDB Setup (Replica Set for Transactions)
1. Stop any running `mongod` process.
2. Start MongoDB as a replica set:
   ```sh
   mongod --replSet rs0 --dbpath "C:\data\db"
   ```
3. In a new terminal:
   ```sh
   mongo
   > rs.initiate()
   ```
4. [Optional] Use MongoDB Compass to view/manage data.

---

## Configuration

Create a `.env` or `config.env` file in the project root. Example:
```
MONGO_URI=mongodb://localhost:27017
MONGO_DB=blogdb
JWT_SECRET=your_jwt_secret
PORT=8080
EMAIL_FROM=your@email.com
EMAIL_PASS=your_email_password
```

---

## Running the App
```sh
go run delivery/main.go
```
The server will start on the port specified in your config (default: 8080).

---

## API Endpoints

### Auth
- `POST /auth/register` — Register a new user
- `POST /auth/login` — Login and receive JWT cookie
- `GET /auth/verify` — Email verification
- `POST /auth/forget` — Request password reset
- `POST /auth/reset` — Reset password
- `POST /auth/logout` — Logout (requires auth)
- `GET /auth/refresh` — Refresh JWT (requires auth)

### Blogs
- `GET /blogs` — List blogs (paginated)
- `GET /blogs/:id` — Get blog by ID (requires auth)
- `POST /blogs` — Create blog (auth: USER/ADMIN)
- `PUT /blogs/:id` — Update blog (auth: USER/ADMIN, must be author)
- `DELETE /blogs/:id` — Delete blog (auth: USER/ADMIN, must be author)
- `GET /blogs/search` — Search blogs by title/author
- `GET /blogs/filter` — Filter blogs by tag, popularity, etc.

### Blog Reactions
- `POST /blogs/:id/like` — Like a blog (auth)
- `POST /blogs/:id/dislike` — Dislike a blog (auth)
- `POST /blogs/:id/unlike` — Remove like (auth)
- `POST /blogs/:id/undislike` — Remove dislike (auth)

### Comments
- `GET /blogs/:id/comments` — List comments for a blog
- `POST /blogs/:id/comments` — Add comment (auth)
- `PUT /comments/:id` — Update comment (auth, must be author)
- `DELETE /blogs/:id/comments/:commentID` — Delete comment (auth, must be author)

### Comment Reactions
- `POST /comments/:id/react/:status` — React to comment (auth)
- `GET /comments/:id/reaction` — Get user reaction (auth)

### Tags
- Tags are auto-created/normalized when creating/updating blogs.

---

## Authentication & Roles
- JWT tokens are issued on login and stored in an `auth_token` cookie.
- Use the cookie in all requests to protected endpoints.
- Roles: `USER`, `ADMIN` (case-sensitive, see domain/user.go)
- Role-based middleware restricts access to certain endpoints.

---

## Testing
- Unit and integration tests are in the `repositories/` and `infrastructures/` folders.
- To run tests:
```sh
go test ./...
```
- Integration tests require a running MongoDB instance.

---

## Development Workflow
- Use feature branches for new features/fixes.
- Pull latest changes: `git pull --rebase`
- Create a new branch: `git checkout -b feature/your-feature`
- Commit and push: `git add . && git commit -m "feat: ..." && git push origin feature/your-feature`
- Open a PR for review.

---

## Troubleshooting
- **Transactions error:** Make sure MongoDB is running as a replica set.
- **User ID not found in context:** Ensure you are sending the `auth_token` cookie.
- **Config file not found:** Create a `.env` or `config.env` in the project root.
- **CORS issues:** Configure Gin CORS middleware as needed.
- **Email not sending:** Check your SMTP credentials in config.

---

## Contributing
1. Fork the repo and clone your fork.
2. Create a new branch for your feature or bugfix.
3. Write clear, well-documented code and tests.
4. Open a pull request with a detailed description.

---

## License
MIT
