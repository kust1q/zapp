# Zapp
A lightweight REST API inspired by X.com, built with GO.

# Technologies
[![Go](https://img.shields.io/badge/-Go-464646?style=flat-square&logo=go)](https://go.dev/)
[![Redis](https://img.shields.io/badge/-Redis-464646?style=flat-square&logo=redis)](https://redis.io/)
[![Kafka](https://img.shields.io/badge/-Kafka-464646?style=flat-square&logo=apache-kafka)](https://kafka.apache.org/)
[![Elasticsearch](https://img.shields.io/badge/-Elasticsearch-464646?style=flat-square&logo=elasticsearch)](https://www.elastic.co/elasticsearch/)
[![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-464646?style=flat-square&logo=postgresql)](https://www.postgresql.org/)
[![MinIO](https://img.shields.io/badge/-MinIO-464646?style=flat-square&logo=minio)](https://min.io/)
[![golang-migrate](https://img.shields.io/badge/-Migrate-464646?style=flat-square&logo=go)](https://github.com/golang-migrate/migrate)
[![Prometheus](https://img.shields.io/badge/-Prometheus-464646?style=flat-square&logo=prometheus)](https://prometheus.io/)
[![Grafana](https://img.shields.io/badge/-Grafana-464646?style=flat-square&logo=grafana)](https://grafana.com/)
[![Gin](https://img.shields.io/badge/-Gin-464646?style=flat-square&logo=go)](https://gin-gonic.com/)
[![gRPC](https://img.shields.io/badge/-gRPC-464646?style=flat-square&logo=grpc)](https://grpc.io/)
[![Govatar](https://img.shields.io/badge/-Govatar-464646?style=flat-square&logo=go)](https://github.com/alexeyco/govatar)
[![Docker](https://img.shields.io/badge/-Docker-464646?style=flat-square&logo=docker)](https://www.docker.com/)
[![Kubernetes](https://img.shields.io/badge/-Kubernetes-464646?style=flat-square&logo=kubernetes)](https://kubernetes.io/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Tech Stack

### Languages & Frameworks
- **Programming Language**: Go 1.24.4
- **Web Framework**: [Gin](https://gin-gonic.com/) — High-performance HTTP framework for building REST APIs
- **gRPC**: Implemented services using the gRPC protocol for inter-service communication
- **Avatar Generation**: [govatar](https://github.com/alexeyco/govatar) — Library for generating random avatars

### Data Storage
- **Relational Database**: [PostgreSQL](https://www.postgresql.org/) — Reliable and scalable RDBMS
- **Caching / Sessions**: [Redis](https://redis.io/) — Fast key-value store for caching and session management
- **Search & Analytics**: [Elasticsearch](https://www.elastic.co/elasticsearch/) — Full-text search and data aggregation engine
- **Object Storage**: [MinIO S3](https://min.io/) — S3-compatible storage for files and media

### Migrations & Schema Management
- **Database Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate) — Tool for managing PostgreSQL schema migrations

### Monitoring & Observability
- **Metrics**: [Prometheus](https://prometheus.io/) — Collection and visualization of application metrics
- **Dashboards**: [Grafana](https://grafana.com/) — Visualization of metrics, logs, and traces

### Event Streaming & Messaging
- **Event Broker**: [Apache Kafka](https://kafka.apache.org/) — Distributed event streaming platform for asynchronous communication between services

### Infrastructure & Deployment
- **Containerization**: [Docker](https://www.docker.com/) — Building and running containers
- **Orchestration**: [Kubernetes](https://kubernetes.io/) — Managing microservices in production

### Configuration
- **Config Management**: [Viper](https://github.com/spf13/viper) — Application configuration via `.env` and YAML config files
- **Logging**: [Logrus](https://github.com/sirupsen/logrus) — Structured logging in JSON format

## Project Structure

```
Zapp/
├── proto/                         # Shared protobuf contracts (source of truth)
│   ├── search/search.proto        #   Search service contract
│   ├── tweet/tweet.proto          #   Tweet service contract
│   └── user/user.proto            #   User service contract
│
├── main/                          # Main microservice
│   ├── cmd/app/
│   │   └── main.go                # Entry point
│   ├── internal/
│   │   ├── config/                # Application configuration
│   │   ├── controllers/           # HTTP handlers, DTOs, converters, middleware
│   │   │   ├── http/              #   REST API controllers
│   │   │   └── grpc/              #   gRPC servers (tweet, user)
│   │   ├── providers/             # Infrastructure providers
│   │   │   ├── db/                #   PostgreSQL, Redis, MinIO
│   │   │   ├── search/            #   gRPC client to search service
│   │   │   └── websocket/         #   WebSocket hub
│   │   ├── service/               # Business logic
│   │   │   ├── auth/              #   Authentication & authorization
│   │   │   ├── tweets/            #   Tweet CRUD, likes, retweets
│   │   │   ├── user/              #   User management, follows
│   │   │   ├── feed/              #   Feed generation
│   │   │   ├── media/             #   Media handling
│   │   │   ├── search/            #   Search orchestration
│   │   │   └── websocket/         #   WebSocket notifications
│   │   ├── domain/
│   │   │   ├── entity/            # Domain entities
│   │   │   └── events/            # Domain events for Kafka
│   │   └── errs/                  # Custom error definitions
│   ├── pkg/                       # Shared infrastructure packages
│   │   ├── kafka/                 # Kafka producer
│   │   ├── postgres/              # PostgreSQL connection helper
│   │   ├── redis/                 # Redis client helper
│   │   ├── minio/                 # MinIO client helper
│   │   ├── elastic/               # Elasticsearch client helper
│   │   └── gen/proto/             # Generated protobuf code
│   ├── migrations/                # SQL migration files
│   ├── configs/                   # YAML config files
│   ├── docs/                      # Swagger documentation
│   ├── Dockerfile
│   └── go.mod
│
├── search/                        # Search microservice
│   ├── cmd/search/
│   │   └── main.go                # Entry point
│   ├── internal/
│   │   ├── config/                # Configuration
│   │   ├── controllers/           # gRPC server, Kafka handler
│   │   │   ├── grpc/              #   SearchService gRPC server
│   │   │   └── kafka/             #   Event consumer handler
│   │   ├── providers/
│   │   │   └── search/
│   │   │       └── elastic/       #   Elasticsearch repository
│   │   ├── service/
│   │   │   └── search/            # Search business logic
│   │   └── domain/
│   │       ├── entity/            # Minimal entitie
│   │       └── events/            # Kafka event definitions
│   ├── pkg/
│   │   ├── elastic/               # Elasticsearch client helper
│   │   ├── kafka/                 # Kafka consumer
│   │   └── gen/proto/             # Generated protobuf code (from /proto)
│   ├── configs/                   # YAML config file
│   ├── Dockerfile
│   └── go.mod
│
├── k8s/                           # Kubernetes manifests
│   ├── namespace.yaml
│   ├── ingress.yaml
│   ├── main/                      #   Main service deployment & configmap
│   ├── search/                    #   Search service deployment & configmap
│   ├── api/                       #   Service & ingress for main
│   ├── postgres/
│   ├── redis/
│   ├── minio/
│   ├── elasticsearch/
│   ├── kafka/
│   ├── prometheus/
│   ├── grafana/
│   └── migrate/
│
├── docker-compose.yml             # Local development
├── .dockerignore
├── .gitignore
└── LICENSE
```

## Quick Start

### Requirements

- Go 1.24.4+
- Docker & Docker Compose
- PostgreSQL 16.11+

### Run with Docker

1. Create `.env` file in `main/`:
```env
POSTGRES_USER=postgres_user
POSTGRES_PASSWORD=postgres_password
POSTGRES_DB=postgres_db
CONTAINER_DB=postgres
DB_URL="postgres://postgres_user:postgres_password@localhost:5432/postgres_db?sslmode=disable"

MINIO_USER=minio_user
MINIO_PASSWORD=minio_password

REDIS_PASSWORD=redis_password

HASH_SECRET=ecf5d5137aa0ce362d8f496c154ff53d31bb1b381a6a740684a987bc5e80647c

PRIVATE_KEY_PATH=./certs/private.pem
PUBLIC_KEY_PATH=./certs/public.pem

KAFKA_CLUSTER_ID=HhCsSRTLRM6R30NVbW5YJQ
```

2. Start services:
```bash
docker-compose up -d --build
```

This starts: PostgreSQL, Redis, MinIO, Elasticsearch, Kafka, the **main** API service, and the **search** microservice.

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/sign-up` | Register new user |
| POST | `/api/v1/auth/sign-in` | Login |
| POST | `/api/v1/auth/sign-out` | Logout |
| POST | `/api/v1/auth/refresh` | Refresh access token |

### Tweets
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tweets` | Create tweet |
| GET | `/api/v1/tweets/:id` | Get tweet by ID |
| GET | `/api/v1/tweets` | Get user's tweets |
| PUT | `/api/v1/tweets/:id` | Update tweet |
| DELETE | `/api/v1/tweets/:id` | Delete tweet |
| POST | `/api/v1/tweets/:id/like` | Like tweet |
| POST | `/api/v1/tweets/:id/retweet` | Retweet |

### User
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/:id` | Get user profile |
| GET | `/api/v1/users/search` | Search users |
| POST | `/api/v1/users/follow/:id` | Follow user |
| POST | `/api/v1/users/unfollow/:id` | Unfollow user |

### Search
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/search/tweets?q=query` | Full-text tweet search |
| GET | `/api/v1/search/users?q=query` | Full-text user search |

## Architecture

The project follows a **clean architecture** approach with two microservices:

### main (REST API + gRPC client)

```
HTTP Client → Gin HTTP Handlers
                  ↓
              Service Layer (auth, tweets, user, feed, search)
                  ↓
        ┌─────────┴──────────┐
        ↓                    ↓
   PostgreSQL           gRPC Client → search service
   Redis Cache          (search API)
   MinIO Storage
```

### search (gRPC server + Kafka consumer)

```
Kafka Events → Kafka Handler
                   ↓
             Search Service (index, delete, search)
                   ↓
              Elasticsearch

gRPC Client → gRPC Server (SearchService)
                  ↓
             Search Service → Elasticsearch
```

### Communication Between Services

1. **Synchronous**: main → search via gRPC for search queries
2. **Asynchronous**: main → search via Kafka for indexing

## License

This project is licensed under the [MIT License](LICENSE).