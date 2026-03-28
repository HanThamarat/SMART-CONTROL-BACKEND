# Smart Control Backend

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

## Overview

This project provides the backend runtime for Smart Control. It uses Fiber for HTTP handling, GORM for database access, PostgreSQL for persistence, and Mosquitto for MQTT communication. The repository already includes a Docker Compose setup for local development, so the full stack can be started with a single command.

## Tech Stack

| Layer | Tooling |
| --- | --- |
| Language | Go 1.25 |
| Web framework | Fiber v2 |
| ORM | GORM |
| Database | PostgreSQL |
| Messaging | Eclipse Mosquitto |
| Containers | Docker Compose |

## Project Structure

```text
.
|-- internal/
|   |-- domain/          # Entities, DTOs, repository/usecase contracts
|   |-- handler/         # Fiber handlers
|   |-- mqttbridge/      # MQTT <-> WebSocket bridge setup
|   |-- repositories/    # GORM-backed repository implementations
|   |-- router/          # Route registration
|   |-- socket/          # Realtime socket server
|   `-- usecase/         # Business logic layer
|-- pkg/
|   |-- database/        # PostgreSQL connection setup
|   |-- encrypt/         # Password hashing and JWT helpers
|   |-- initialize/      # Default data initialization
|   |-- loadEnd/         # .env loader
|   |-- mqttCon/         # MQTT connection helper
|   `-- response/        # Shared response formatter
|-- Dockerfile.dev       # Development image with Air
|-- Dockerfile.prod      # Production image
|-- docker-compose.yaml  # Local stack definition
|-- main.go              # Application bootstrap
`-- .env.example         # Example environment file
```

## Environment Setup

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

MQTT_BROKER=mosquitto:11883
MQTT_CLIENT_ID=smart-control-backend
MQTT_USERNAME=smartcontrol
MQTT_PASSWORD=123456
MQTT_SUBSCRIBE_TOPICS=TEST/MQTT

EMAIL=admin@gmail.com
NAME=Administrator
USERNAME=admin
PASSWORD=12345678
```

Important notes:

- `Port` should remain `:3000` to match the compose port mapping.
- `DATABASE_HOST` should remain `postgreSql` when the app runs inside Docker.
- On first startup, the app auto-migrates the database and seeds a default user from the env values above.

## Quick Start

### 1. Prepare the environment

```bash
cp .env.example .env
```

### 2. Start the full stack

```bash
docker compose up --build
```

### 3. Available services

| Service | Address |
| --- | --- |
| Backend | `http://localhost:3000` |
| PostgreSQL | `localhost:5432` |
| MQTT | `localhost:1883` |
| MQTT WebSocket | `localhost:9001` |

### 4. Verify the backend is running

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

## Docker Compose Services

The local stack includes three services:

- `server` for the Go backend with hot reload via Air
- `postgreSql` for relational data storage
- `mosquitto` for MQTT messaging

The backend service waits for the database and broker health checks before starting.

## Development Notes

- `Dockerfile.dev` is the default image used by `docker-compose.yaml`.
- `Dockerfile.prod` can be used by setting `DOCKERFILE=Dockerfile.prod`.
- The application loads configuration from `.env` during startup.

## Frontend WebSocket to MQTT

The backend exposes a WebSocket endpoint at `/socket.io` and forwards socket messages to MQTT.

Authentication:

- The socket endpoint requires a JWT.
- Pass the token as `?token=YOUR_JWT` in the WebSocket URL.

Example connection:

```js
const ws = new WebSocket("ws://localhost:3000/socket.io?token=YOUR_JWT");

ws.onopen = () => {
  console.log("socket connected");
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log("socket message:", message);
};
```

### Publish to MQTT from frontend

Send a WebSocket message with event `mqtt:publish`:

```js
ws.send(JSON.stringify({
  event: "mqtt:publish",
  data: {
    topic: "TEST/MQTT",
    payload: "ON",
    qos: 0,
    retained: false
  }
}));
```

Example with JSON payload:

```js
ws.send(JSON.stringify({
  event: "mqtt:publish",
  data: {
    topic: "TEST/MQTT",
    payload: {
      led: true,
      value: 1
    }
  }
}));
```

If publish succeeds, the backend responds with:

```json
{
  "event": "mqtt:published",
  "message": "MQTT message published."
}
```

### Receive MQTT messages in frontend

When the backend receives an MQTT message from a subscribed topic, it forwards it to connected socket clients as `mqtt:message`.

Example received message:

```json
{
  "event": "mqtt:message",
  "message": "MQTT message received.",
  "timestamp": "2026-03-28T00:00:00Z",
  "data": {
    "topic": "TEST/MQTT",
    "payload": "ON"
  }
}
```

Example listener:

```js
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);

  if (message.event === "mqtt:message") {
    console.log("MQTT topic:", message.data.topic);
    console.log("MQTT payload:", message.data.payload);
  }
};
```

### MQTT topics

- Configure subscribed topics with `MQTT_SUBSCRIBE_TOPICS`.
- Use comma-separated values for multiple topics, for example `TEST/MQTT,smart/control,sensor/temp`.
- If not set, the backend subscribes to `TEST/MQTT`.

## Stop the Stack

```bash
docker compose down
```

Remove containers and volumes:

```bash
docker compose down -v
```

## GitHub Actions CI/CD

This repository now includes a GitHub Actions workflow at `.github/workflows/ci-cd.yml`.

What it does:

- Runs CI on `pull_request`, `main`, and `master`
- Checks Go formatting with `gofmt`
- Runs `go test ./...`
- Deploys to your Ubuntu server on pushes to `main` or `master`

### Deployment flow

The deploy job:

- connects to your Ubuntu server over SSH
- syncs the repository to your server
- writes a production `.env` file from a GitHub Secret
- runs `docker compose -f docker-compose.prod.yaml up -d --build`

### Required GitHub Secrets

Add these in `GitHub -> Settings -> Secrets and variables -> Actions`:

- `SERVER_HOST` = your Ubuntu server IP or domain
- `SERVER_PORT` = your SSH port, usually `22`
- `SERVER_USER` = the Linux user used for deployment
- `SERVER_SSH_KEY` = the private SSH key used by GitHub Actions
- `SERVER_APP_DIR` = target app directory on the server, for example `/home/ubuntu/apps/smart-control-backend`
- `ENV_FILE` = the full content of your production `.env`

Example `ENV_FILE` secret:

```env
Port=:3000
SERVER_MODE=false
JWT_SECRET=change-me

DATABASE_HOST=postgreSql
DATABASE_USERNAME=root
DATABASE_PASSWORD=root
DATABASE_NAME=smartcontrol
DATABASE_PORT=5432

MQTT_BROKER=mosquitto:1883
MQTT_CLIENT_ID=smart-control-backend
MQTT_USERNAME=smartcontrol
MQTT_PASSWORD=123456
MQTT_SUBSCRIBE_TOPICS=TEST/MQTT

EMAIL=admin@example.com
NAME=Admin
USERNAME=admin
PASSWORD=strong-password
```

### Ubuntu server requirements

Your server should already have:

- Docker installed
- Docker Compose plugin installed
- the deploy user able to run `docker compose`

Manual first-time setup example:

```bash
mkdir -p /home/ubuntu/apps/smart-control-backend
```
