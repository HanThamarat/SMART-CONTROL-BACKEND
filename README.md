# Smart Control Backend 🚀

<p align="center">
  Backend service for Smart Control, built with Go, Fiber, PostgreSQL, and MQTT.
</p>

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="Fiber" src="https://img.shields.io/badge/Fiber-v2-1D7AFC?style=for-the-badge" />
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-18-336791?style=for-the-badge&logo=postgresql&logoColor=white" />
  <img alt="Mosquitto" src="https://img.shields.io/badge/MQTT-Mosquitto-3C5280?style=for-the-badge" />
  <img alt="Docker Compose" src="https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
</p>

## Overview 🌐

This project provides the backend runtime for Smart Control. It uses Fiber for HTTP handling, GORM for database access, PostgreSQL for persistence, and Mosquitto for MQTT communication. The repository already includes a Docker Compose setup for local development, so the full stack can be started with a single command.

## Tech Stack 🧰

| Layer | Tooling |
| --- | --- |
| Language | Go 1.25 |
| Web framework | Fiber v2 |
| ORM | GORM |
| Database | PostgreSQL |
| Messaging | Eclipse Mosquitto |
| Containers | Docker Compose |

## Project Structure 🗂️

```text
.
├── internal/
│   ├── domain/          # Entities, DTOs, repository/usecase contracts
│   ├── handler/         # Fiber handlers
│   ├── repositories/    # GORM-backed repository implementations
│   ├── router/          # Route registration
│   └── usecase/         # Business logic layer
├── pkg/
│   ├── database/        # PostgreSQL connection setup
│   ├── encrypt/         # Password hashing and JWT helpers
│   ├── initialize/      # Default data initialization
│   ├── loadEnd/         # .env loader
│   ├── mqttCon/         # MQTT connection helper
│   └── response/        # Shared response formatter
├── Dockerfile.dev       # Development image with Air
├── Dockerfile.prod      # Production image
├── docker-compose.yaml  # Local stack definition
├── main.go              # Application bootstrap
└── .env.example         # Example environment file
```

## Environment Setup ⚙️

Create your local environment file:

```bash
cp .env.example .env
```

Recommended values for Docker Compose:

```env
Port=:3000
SERVER_MODE=false
JWT_SECRET=JWT_SECRET

DATABASE_HOST=postgreSql
DATABASE_USERNAME=root
DATABASE_PASSWORD=root
DATABASE_NAME=smartcontrol
DATABASE_PORT=5432

MQTT_BROKER=mosquitto:1883
MQTT_CLIENT_ID=smart-control-backend
MQTT_USERNAME=smartcontrol
MQTT_PASSWORD=123456

EMAIL=HanThamarat@gmail.com
NAME=HanThamarat
USERNAME=HanThamarat
PASSWORD=123456
```

Important notes:

- `Port` should remain `:3000` to match the compose port mapping.
- `DATABASE_HOST` should remain `postgreSql` when the app runs inside Docker.
- On first startup, the app auto-migrates the database and seeds a default user from the env values above.

## Quick Start ⚡

### 1. Prepare the environment 📄

```bash
cp .env.example .env
```

### 2. Start the full stack 🐳

```bash
docker compose up --build
```

### 3. Available services 📡

| Service | Address |
| --- | --- |
| Backend | `http://localhost:3000` |
| PostgreSQL | `localhost:5432` |
| MQTT | `localhost:1883` |
| MQTT WebSocket | `localhost:9001` |

### 4. Verify the backend is running ✅

```bash
curl http://localhost:3000/
```

Expected response:

```json
{
  "status": 200,
  "message": "Server is running.",
  "body": null
}
```

## Docker Compose Services 🧱

The local stack includes three services:

- `server` for the Go backend with hot reload via Air
- `postgreSql` for relational data storage
- `mosquitto` for MQTT messaging

The backend service waits for the database and broker health checks before starting.

## Development Notes 🛠️

- `Dockerfile.dev` is the default image used by `docker-compose.yaml`.
- `Dockerfile.prod` can be used by setting `DOCKERFILE=Dockerfile.prod`.
- The application loads configuration from `.env` during startup.

## Stop the Stack 🛑

```bash
docker compose down
```

Remove containers and volumes:

```bash
docker compose down -v
```
