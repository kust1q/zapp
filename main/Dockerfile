FROM golang:1.24.4 AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc-dev \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o zapp-app ./cmd/app

FROM alpine:3.21

RUN apk add --no-cache \
    ca-certificates \
    libpq \
    libc6-compat \
    && rm -rf /var/cache/apk/*

RUN mkdir -p /app/certs /app/configs

COPY --from=builder /app/zapp-app /app/zapp-app
COPY --from=builder /app/certs /app/certs
COPY --from=builder /app/configs /app/configs

RUN chmod 755 /app/zapp-app && chmod 644 /app/certs/*.pem /app/configs/*.yml

WORKDIR /app

EXPOSE 8080 50051

CMD ["/app/zapp-app"]
