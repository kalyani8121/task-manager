FROM golang:1.25.4-alpine AS builder

# Install git (needed for go mod download)
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files FIRST (before source code).
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/task-manager ./cmd/main.go


FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/task-manager .

COPY migrations/ ./migrations/

# Expose the port our app listens on
EXPOSE 8080

# The command to run when the container starts
CMD ["./task-manager"]