FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o zapp-search ./cmd/search

FROM alpine:3.21

RUN apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app/configs

COPY --from=builder /app/zapp-search /app/zapp-search
COPY --from=builder /app/configs /app/configs

RUN chmod 755 /app/zapp-search && chmod 644 /app/configs/*.yml

WORKDIR /app

EXPOSE 50052

CMD ["/app/zapp-search"]
