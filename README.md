# Subscription Service

A microservice for managing subscriptions, billing, and usage tracking in the ERP Suite.

## 🏗️ Architecture

The subscription service follows Clean Architecture principles with a layered approach:

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ HTTP/REST   │  │ gRPC Server │  │ GraphQL Resolver    │ │
│  │ Handlers    │  │             │  │                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ Plan        │  │ Subscription│  │ Billing             │ │
│  │ Service     │  │ Service     │  │ Service             │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ Entities    │  │ Value       │  │ Domain              │ │
│  │             │  │ Objects     │  │ Services            │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ Database    │  │ Cache       │  │ External            │ │
│  │ Repositories│  │ (Redis)     │  │ Services            │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
subscription-service/
├── cmd/                           # Application entry points
│   └── subscription-service/      # Main application
├── internal/                      # Private application code
│   ├── domain/                    # Domain layer (business logic)
│   │   ├── entities/              # Domain entities
│   │   ├── interfaces/            # Domain interfaces
│   │   └── services/              # Domain services
│   ├── application/               # Application layer (use cases)
│   │   ├── dto/                   # Data Transfer Objects
│   │   ├── interfaces/            # Application interfaces
│   │   └── services/              # Application services
│   └── infrastructure/            # Infrastructure layer
│       ├── interfaces/            # Infrastructure interfaces
│       └── repositories/          # Data access implementations
├── pkg/                          # Public packages
│   ├── config/                   # Configuration management
│   ├── database/                 # Database connection
│   ├── redis/                    # Redis client
│   ├── grpc/                     # gRPC server
│   ├── handlers/                 # HTTP handlers
│   ├── middleware/               # HTTP middleware
│   └── utils/                    # Utility functions
├── proto/                        # Protocol buffer definitions
├── migrations/                   # Database migrations
├── scripts/                      # Utility scripts
├── tests/                        # Integration tests
├── main.go                       # Application entry point
├── Makefile                      # Build and development commands
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── Dockerfile                    # Docker configuration
├── docker-compose.yml            # Docker Compose configuration
├── .env.example                  # Environment variables example
└── README.md                     # This file
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher
- Protocol Buffers compiler (protoc)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd subscription-service
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Set up environment**
   ```bash
   make env-setup
   # Edit .env file with your configuration
   ```

4. **Set up database**
   ```bash
   make db-create
   make migrate
   make seed
   ```

5. **Generate protobuf files**
   ```bash
   make proto
   ```

6. **Run the service**
   ```bash
   make run
   ```

### Development

```bash
# Run with hot reload
make dev

# Run tests
make test

# Run tests with coverage
make test-cover

# Format code
make fmt

# Lint code
make lint
```

## 📊 Features

### Subscription Management
- Create, update, and cancel subscriptions
- Support for different billing cycles (monthly, yearly)
- Trial period management
- Plan changes with proration
- Subscription status tracking

### Plan Management
- Create and manage subscription plans
- Feature-based plan configuration
- Pricing management
- Plan availability controls

### Billing & Invoicing
- Automatic invoice generation
- Payment processing
- Invoice status tracking
- Payment method management
- Tax and discount support

### Usage Tracking
- Real-time usage monitoring
- Usage limits enforcement
- Usage history and analytics
- Multiple metric types (users, storage, API calls)

### Payment Methods
- Credit card support
- PayPal integration
- Payment method validation
- Default payment method management

## 🔧 Configuration

The service can be configured using environment variables or a configuration file:

### Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=erp_subscription

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# HTTP Server
HTTP_PORT=8081

# gRPC Server
GRPC_PORT=50051

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=3600

# Environment
ENVIRONMENT=development
```

### Configuration File

Create a `config.yaml` file:

```yaml
environment: development

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: erp_subscription
  sslmode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

http:
  port: 8081

grpc:
  port: 50051

jwt:
  secret: your-secret-key
  expiration: 3600

log:
  level: info
```

## 🗄️ Database Schema

### Core Tables

- **plans**: Subscription plans with features and pricing
- **subscriptions**: Organization subscriptions
- **invoices**: Billing invoices
- **payment_methods**: Payment methods for organizations
- **usage**: Usage tracking data

### Key Relationships

- Organizations have one active subscription
- Subscriptions belong to a plan
- Invoices are generated for subscriptions
- Payment methods belong to organizations
- Usage is tracked per subscription

## 🔌 API Endpoints

### HTTP REST API

#### Health Check
```
GET /health
```

#### Plans
```
GET    /api/v1/plans
GET    /api/v1/plans/:id
POST   /api/v1/plans
PUT    /api/v1/plans/:id
DELETE /api/v1/plans/:id
```

#### Subscriptions
```
GET    /api/v1/subscriptions/:organization_id
POST   /api/v1/subscriptions
PUT    /api/v1/subscriptions/:id
DELETE /api/v1/subscriptions/:id
POST   /api/v1/subscriptions/:id/cancel
POST   /api/v1/subscriptions/:id/reactivate
POST   /api/v1/subscriptions/:id/change-plan
```

#### Invoices
```
GET    /api/v1/invoices
GET    /api/v1/invoices/:id
POST   /api/v1/invoices
POST   /api/v1/invoices/:id/pay
```

#### Payment Methods
```
GET    /api/v1/payment-methods
GET    /api/v1/payment-methods/:id
POST   /api/v1/payment-methods
PUT    /api/v1/payment-methods/:id
DELETE /api/v1/payment-methods/:id
POST   /api/v1/payment-methods/:id/set-default
```

#### Usage
```
GET    /api/v1/usage
POST   /api/v1/usage/track
GET    /api/v1/usage/history
```

### gRPC API

The service also exposes a gRPC API with the same functionality. See `proto/subscription.proto` for the complete API definition.

## 🧪 Testing

### Unit Tests
```bash
make test
```

### Integration Tests
```bash
make test-integration
```

### Coverage Report
```bash
make test-cover
```

## 🐳 Docker

### Build Image
```bash
make docker-build
```

### Run Container
```bash
make docker-run
```

### Docker Compose
```bash
docker-compose up -d
```

## 📈 Monitoring

### Health Checks
- HTTP: `GET /health`
- gRPC: `HealthCheck` method

### Metrics
- Prometheus metrics available at `/metrics`
- Custom business metrics for subscriptions, billing, and usage

### Logging
- Structured JSON logging
- Configurable log levels
- Request tracing with correlation IDs

## 🔒 Security

- JWT-based authentication
- Role-based access control
- Input validation and sanitization
- SQL injection prevention
- Rate limiting
- CORS configuration

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run the test suite
6. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the test cases for usage examples 