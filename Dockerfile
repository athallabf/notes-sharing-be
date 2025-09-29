# --- Build Stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy only dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY src/ ./src/
COPY test/ ./test/

# Build the application
# We build from src/main.go, so the output path is just 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./src/main.go


# --- Final Stage ---
FROM alpine:3.20

# Add a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install curl for healthcheck and any other necessary certs
RUN apk add --no-cache curl ca-certificates

WORKDIR /home/appuser

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/main .

# Change ownership of the app directory
USER appuser

# Expose the port the app runs on (should match APP_PORT in .env)
EXPOSE 3000

# The command to run when the container starts.
# Environment variables will be injected by Docker Compose.
CMD ["./main"]