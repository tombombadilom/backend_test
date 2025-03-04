# Backend Project Analysis: Golang Live Ops Events System

## 📑 Navigation

- 🏠 **[Project Home](../../README.md)** - Return to project root
- 📘 **[Documentation Home](../README.md)** - Main documentation
- 📋 **[Instructions](../0-Instructions/INSTRUCTIONS_Backend.md)** - Project requirements
- 🔄 **[Methodology](../1-Methodology/README.md)** - Development approach
- 🏗️ **[Architecture](../1-Design/Architecture.md)** - System design
- 💾 **[Data Models](../1-Design/DataModels.md)** - Data structures
- ⚙️ **[Technical Specs](../3-Specifications/TechnicalSpecifications.md)** - Technical details
- 🔌 **[API Specs](../3-Specifications/APISpecifications.md)** - API documentation
- 🔒 **[Security](../3-Specifications/SecuritySpecifications.md)** - Security measures
- 📝 **[Implementation Plan](../4-Todo/README.md)** - Development tasks

---

## Project Overview

This project requires the development of a Golang backend application that serves both HTTP and gRPC APIs on the same port while implementing role-based access control using API keys. The application is designed to manage player Live Ops events for video games, which are used to dynamically serve limited-time events, offers, and challenges.

## Core Requirements Analysis

### Technical Stack
- **Language**: Golang
- **Database**: SQLite
- **API Protocols**: HTTP REST and gRPC (on the same port)
- **Authentication**: API Key-based with role-based access control
- **Data Format**: Protobuf for gRPC definitions

### Key Functionalities

1. **Dual API Support**:
   - HTTP API for public access (retrieving live events)
   - gRPC API for internal administrative operations (CRUD operations)
   - Both APIs must run on the same port

2. **Role-Based Access Control**:
   - API keys with specific roles:
     - `http_user`: Can only access HTTP endpoints
     - `grpc_admin`: Can only access gRPC endpoints
   - Authentication via `Authorization` header

3. **Data Management**:
   - Persistent storage using SQLite
   - Live Events entity with fields:
     - `id` (UUID)
     - `title` (string)
     - `description` (string)
     - `start_time` (timestamp)
     - `end_time` (timestamp)
     - `rewards` (JSON string)

### API Specifications

#### HTTP Endpoints (Public)
- `GET /events`: Retrieve all active events
- `GET /events/{id}`: Retrieve details of a specific event

#### gRPC Endpoints (Internal)
- `CreateEvent`: Create a new live event
- `UpdateEvent`: Update an existing event
- `DeleteEvent`: Delete an event
- `ListEvents`: Retrieve all events (including past events)

## Technical Challenges

1. **Serving HTTP and gRPC on the same port**:
   - This requires careful implementation of a multiplexer that can distinguish between HTTP and gRPC traffic
   - Need to handle protocol detection and routing appropriately

2. **Role-Based Access Control**:
   - Implementing a secure authentication system using API keys
   - Ensuring proper validation and authorization based on roles
   - Protecting against unauthorized access

3. **Concurrent Request Handling**:
   - Efficiently managing both HTTP and gRPC requests simultaneously
   - Ensuring thread safety and proper resource management

4. **Data Persistence**:
   - Designing an efficient SQLite schema for the Live Events data
   - Implementing proper database operations and error handling

## Implementation Strategy

1. **Architecture Design**:
   - Define clear separation of concerns (controllers, services, repositories)
   - Design a modular system that allows for easy extension and maintenance

2. **API Implementation**:
   - Define Protobuf schemas for the gRPC API
   - Implement HTTP handlers for the REST API
   - Create a multiplexer to handle both protocols on the same port

3. **Authentication System**:
   - Implement API key validation and role-based access control
   - Create middleware for authentication and authorization

4. **Database Layer**:
   - Design and implement SQLite schema
   - Create repository layer for data access

5. **Testing Strategy**:
   - Unit tests for core business logic
   - Integration tests for API endpoints
   - Authentication and authorization tests

## Deliverables

1. **Golang Codebase**:
   - Well-structured and maintainable code
   - Proper separation of concerns
   - Clean error handling and logging

2. **Documentation**:
   - README.md with build, run, and test instructions
   - API documentation for both HTTP and gRPC
   - Authentication details and API key role information

3. **Containerization**:
   - Dockerfile for easy deployment

4. **Tests**:
   - Comprehensive test suite covering core functionality

## Evaluation Metrics

The project will be evaluated based on:
- Code structure and adherence to Go best practices
- Concurrency handling and efficiency
- API design and logical structure
- Security implementation (API key authentication)
- Database handling and SQLite usage
- Documentation quality and completeness

## Next Steps

1. Set up the project structure and basic architecture
2. Define the Protobuf schemas for the gRPC API
3. Implement the database layer with SQLite
4. Create the authentication system with API key validation
5. Implement the HTTP and gRPC servers with proper multiplexing
6. Add comprehensive tests and documentation
7. Create the Dockerfile for containerization 