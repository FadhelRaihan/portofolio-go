# =========================
# Stage 1: Builder
# =========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# =========================
# Stage 2: Runtime
# =========================
FROM alpine:3.19

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/main .

ENV PORT=3000
EXPOSE 3000

CMD ["./main"]
