# =========================
# Build stage
# =========================
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o myapp ./cmd/api

# =========================
# Run stage
# =========================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/myapp /app/myapp
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nonroot:nonroot

EXPOSE 8000

ENTRYPOINT ["./myapp"]
