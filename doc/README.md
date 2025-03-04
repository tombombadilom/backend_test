# Live Ops Events System Documentation

## 📑 Navigation

- 🏠 **[Project Home](../README.md)** - Return to project root
- 📋 **[Instructions](0-Instructions/INSTRUCTIONS_Backend.md)** - Project requirements
- 🔄 **[Methodology](1-Methodology/README.md)** - Development approach
- 🏗️ **[Architecture](1-Design/Architecture.md)** - System design
- 💾 **[Data Models](1-Design/DataModels.md)** - Data structures
- 📊 **[Analysis](2-Analysis/Analysis.md)** - Requirements analysis
- ⚙️ **[Technical Specs](3-Specifications/TechnicalSpecifications.md)** - Technical details
- 🔌 **[API Specs](3-Specifications/APISpecifications.md)** - API documentation
- 🔒 **[Security](3-Specifications/SecuritySpecifications.md)** - Security measures
- 📝 **[Todo List](4-Todo/README.md)** - Development tasks

---

## Project Overview

This documentation presents the analysis, design, and specifications for a Golang backend application that serves both HTTP and gRPC APIs on the same port while enforcing role-based access control using API keys. The application is designed to manage player Live Ops events for video games, which are used to dynamically serve limited-time events, offers, and challenges.

## Test Context

This project is a technical test to be completed within a 2-day timeframe. The implementation focuses exclusively on the core requirements specified in the instructions, prioritizing essential functionality over additional features.

## Documentation Structure

This documentation is organized into the following sections:

### [0. Instructions](0-Instructions/INSTRUCTIONS_Backend.md)
- Project requirements and specifications
- Evaluation criteria
- Submission guidelines

### [1. Methodology](1-Methodology/README.md)
- Development approach
- Project management
- Quality assurance

### [1. Design](1-Design/Architecture.md)
- System architecture
- Data models
- API design

### [2. Analysis](2-Analysis/Analysis.md)
- Requirements analysis
- Technical challenges
- Implementation strategy

### [3. Specifications](3-Specifications/TechnicalSpecifications.md)
- Technical specifications
- API documentation
- Authentication details
- [Security specifications](3-Specifications/SecuritySpecifications.md)

### [4. Todo](4-Todo/README.md)
- Implementation tasks
- Testing tasks
- Documentation tasks

## Key Points

### Technical Stack
- **Language**: Golang
- **Database**: SQLite
- **API Protocols**: HTTP REST and gRPC (on the same port)
- **Authentication**: API Key-based with role-based access control
- **Data Format**: Protobuf for gRPC definitions

### Key Features
- Dual API support (HTTP and gRPC) on the same port
- Role-based access control using API keys
- Live Events management for video games
- Persistent storage using SQLite

### Implementation Priorities
1. Set up the project structure and basic architecture
2. Define the Protobuf schemas for the gRPC API
3. Implement the database layer with SQLite
4. Create the authentication system with API key validation
5. Implement the HTTP and gRPC servers with proper multiplexing
6. Add essential tests and documentation
7. Create the Dockerfile for containerization

## Implementation Timeline

The implementation is structured to be completed in 2 days:

- **Day 1**: Focus on core functionality and API implementation
- **Day 2**: Focus on authentication, testing, and documentation

## Next Steps

1. Review the detailed [Architecture Design](1-Design/Architecture.md)
2. Explore the [API Specifications](3-Specifications/APISpecifications.md)
3. Check the [Implementation Plan](4-Todo/README.md)
4. Review the [Security Specifications](3-Specifications/SecuritySpecifications.md) 