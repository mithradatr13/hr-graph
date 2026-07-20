```markdown
# Task Manager API

A production-ready, clean-architecture Task Manager API built with Go, Gin, GORM (PostgreSQL), and Redis caching. Designed with strict Test-Driven Development (TDD) principles, robust error handling, Prometheus metrics, and live pprof profiling support.

---

## 🛠️ Tech Stack

* **Language:** Go 1.23+
* **Framework:** Gin (HTTP Router & Web Framework)
* **Database & ORM:** PostgreSQL with GORM
* **Caching:** Redis
* **Observability:** Prometheus Metrics & Go pprof profiling
* **Testing:** Testify (`assert`, `mock`) & `sqlmock`
* **Containerization:** Docker & Docker Compose

---

## 📁 Project Architecture

The project strictly follows **Clean Architecture** principles, maintaining a clear separation of concerns across layers:

```text
task-manager/
├── cmd/
│   └── api/                # Application entrypoint
├── internal/
│   ├── domain/             # Business models and interface definitions
│   ├── repository/         # Database implementation (GORM + PostgreSQL)
│   ├── service/            # Business logic layer
│   ├── handler/            # HTTP transport layer (Gin handlers)
│   ├── router/             # Route configurations, Swagger UI, metrics, and pprof
│   └── middleware/         # Custom middlewares (e.g., metrics tracking)
├── pkg/
│   ├── cache/              # Redis cache client implementation
│   ├── config/             # Environment configuration loader
│   └── database/           # PostgreSQL connection initializer
├── docs/
│   └── openapi.yaml        # OpenAPI 3.0 specification file
├── Dockerfile              # Multi-stage Docker build file
├── docker-compose.yml      # Service orchestration
└── go.mod

```

---

## 🚀 Getting Started

### Prerequisites

Ensure you have the following installed on your machine:

* **Docker** & **Docker Compose** (Recommended for containerized deployment)
* **Go 1.23+** (If running locally without Docker)

---

### Running with Docker Compose (Recommended)

1. Clone the repository and navigate to the project root.
2. Build and start all services (API, PostgreSQL, Redis) in detached mode:
```bash
docker compose up --build -d

```


3. To restart the services after making code or configuration changes:
```bash
docker compose up --build -d

```


4. To shut down the services completely:
```bash
docker compose down

```



---

### Running Locally

1. Set up your local PostgreSQL and Redis instances, then configure your environment variables.
2. Run database migrations or start the application entrypoint:
```bash
go run cmd/api/main.go

```



---

## 🧪 Running Tests & Coverage

The project maintains high test coverage through rigorous unit, repository, and handler tests.

To run the entire test suite and generate a coverage report, execute:

```bash
go test -v -coverprofile=coverage.out ./...

```

To view a detailed breakdown of statement coverage per package:

```bash
go tool cover -func=coverage.out

```

---

## 📚 API Documentation (Swagger UI)

Once the application is running, interactive API documentation is automatically served via Swagger UI.

* **Swagger UI Dashboard:** [http://localhost:8080/docs](http://localhost:8080/docs)
* **OpenAPI Raw Spec:** [http://localhost:8080/docs/openapi.yaml](http://localhost:8080/docs/openapi.yaml)

---

## 📊 Monitoring & Profiling

* **Prometheus Metrics:** Available for scraping at `http://localhost:8080/metrics`
* **Profiling (pprof):** Available at `http://localhost:8080/debug/pprof/` for performance analysis and load testing diagnostics.

```

```