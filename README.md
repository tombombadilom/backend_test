# Live Ops Events System

## 📑 Navigation

- 📘 **[Documentation](doc/README.md)** - Comprehensive project documentation
- 🏗️ **[Architecture](doc/1-Design/Architecture.md)** - System design and components
- 💾 **[Data Models](doc/1-Design/DataModels.md)** - Database schema and data structures
- 🔌 **[API Specifications](doc/3-Specifications/APISpecifications.md)** - HTTP and gRPC API details
- 🔒 **[Security](doc/3-Specifications/SecuritySpecifications.md)** - Authentication and security measures
- 📝 **[Todo List](doc/4-Todo/README.md)** - Development roadmap and tasks

---

A Golang backend application that serves both HTTP and gRPC APIs on the same port while enforcing role-based access control using API keys. The application manages player Live Ops events for video games, which are used to dynamically serve limited-time events, offers, and challenges.

## Features

- **Dual API Support**: HTTP API for public access and gRPC API for internal administrative operations, both running on the same port.
- **Role-Based Access Control**: API keys with specific roles for HTTP and gRPC access.
- **Persistent Storage**: SQLite database for storing live events and API keys.
- **Containerization**: Docker support for easy deployment.

## Getting Started

### Prerequisites

- Go 1.21 or later
- SQLite
- Docker (optional)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/liveops.git
cd liveops
```

2. Install dependencies:

```bash
go mod download
```

3. Build the application:

```bash
go build -o liveops ./cmd/server
```

### Running the Application

1. Run the application:

```bash
./liveops
```

2. The application will start and listen on port 8080 by default.

### Configuration

The application can be configured using environment variables:

- `PORT`: The port to listen on (default: 8080)
- `DB_PATH`: The path to the SQLite database file (default: ./data/liveops.db)
- `LOG_LEVEL`: The log level (default: info)
- `API_KEY_HTTP`: The API key for HTTP access (for development only)
- `API_KEY_GRPC`: The API key for gRPC access (for development only)

### Docker

1. Build the Docker image:

```bash
docker build -t liveops .
```

2. Run the Docker container:

```bash
docker run -p 8080:8080 liveops
```

## API Documentation

### HTTP API

#### GET /events

Retrieves all active events.

**Authentication**: Requires API Key with `http_user` role.

**Request**:
```http
GET /events HTTP/1.1
Host: example.com
Authorization: Bearer <api_key>
```

**Response**:
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "events": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Summer Festival",
      "description": "Join the summer festival and earn exclusive rewards!",
      "start_time": "2023-06-01T00:00:00Z",
      "end_time": "2023-06-30T23:59:59Z",
      "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}"
    }
  ]
}
```

#### GET /events/{id}

Retrieves details of a specific event.

**Authentication**: Requires API Key with `http_user` role.

**Request**:
```http
GET /events/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
Host: example.com
Authorization: Bearer <api_key>
```

**Response**:
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Summer Festival",
  "description": "Join the summer festival and earn exclusive rewards!",
  "start_time": "2023-06-01T00:00:00Z",
  "end_time": "2023-06-30T23:59:59Z",
  "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}"
}
```

### gRPC API

The gRPC API provides methods for creating, updating, and deleting live events. It requires an API key with the `grpc_admin` role.

For detailed information about the gRPC API, please refer to the [API Specifications](doc/3-Specifications/APISpecifications.md).

## Authentication

All API requests must include an API key in the `Authorization` header. The format of the header is:

```
Authorization: Bearer <api_key>
```

API keys have roles:
- `http_user`: Can access only HTTP endpoints.
- `grpc_admin`: Can access only gRPC endpoints.

## Development

### Project Structure

The project follows a standard Go project layout:

```
/
├── cmd/                  # Main applications
│   └── server/           # The server application
├── internal/             # Private application and library code
│   ├── api/              # API handlers (HTTP and gRPC)
│   ├── auth/             # Authentication and authorization
│   ├── config/           # Configuration
│   ├── db/               # Database access
│   ├── models/           # Domain models
│   └── service/          # Business logic
├── pkg/                  # Public library code
│   └── proto/            # Protobuf definitions
├── scripts/              # Scripts for development and CI/CD
├── test/                 # Additional test applications and test data
└── vendor/               # Application dependencies
```

### Testing

Run the tests:

```bash
go test ./...
```

## Documentation

For detailed documentation about the project, please refer to the [Documentation](doc/README.md).

## License

This project is licensed under the MIT License - see the LICENSE file for details.
