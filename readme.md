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

* **Prometheus Metrics:** Available for scraping at `http://localhost:9090/metrics`
* **Profiling (pprof):** Available at `http://localhost:8080/debug/pprof/` for performance analysis and load testing diagnostics.
