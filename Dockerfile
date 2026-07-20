# Stage 1: Build the binary
FROM docker.arvancloud.ir/library/golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Stage 2: Run the binary in a clean scratch environment
FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
EXPOSE 8080
CMD ["./main"]