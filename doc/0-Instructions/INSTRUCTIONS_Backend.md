# **Technical Test Statement: Golang Backend Developer**

## 📑 Navigation

- 🏠 **[Project Home](../../README.md)** - Return to project root
- 📘 **[Documentation Home](../README.md)** - Main documentation
- 🔄 **[Methodology](../1-Methodology/README.md)** - Development approach
- 🏗️ **[Architecture](../1-Design/Architecture.md)** - System design
- 💾 **[Data Models](../1-Design/DataModels.md)** - Data structures
- 📊 **[Analysis](../2-Analysis/Analysis.md)** - Requirements analysis
- ⚙️ **[Technical Specs](../3-Specifications/TechnicalSpecifications.md)** - Technical details
- 🔌 **[API Specs](../3-Specifications/APISpecifications.md)** - API documentation
- 🔒 **[Security](../3-Specifications/SecuritySpecifications.md)** - Security measures
- 📝 **[Implementation Plan](../4-Todo/README.md)** - Development tasks

---

## **Objective**

Design and develop a **Golang backend application** that serves **both HTTP and gRPC APIs on the same port** while
enforcing **role-based access control using API keys**. The application will manage **player Live Ops events**, used to
dynamically serve limited-time events, offers, and challenges in a video game.

---

## **Requirements**

You are required to build a **Golang application** that:

1. **Exposes an HTTP API** to **retrieve live events for players**.
2. **Exposes a gRPC API** to **create, update, and delete live events** (internal usage).
3. **Runs both HTTP and gRPC servers on the same port**.
4. **Uses API keys** to distinguish between **public users (HTTP access)** and **internal users (gRPC access)**.
5. **Persists data using SQLite**.
6. **Uses Protobuf definitions** for gRPC.
7. **Includes unit tests** where relevant.

---

## **Specifications: Live Ops Events System**

The system should allow managing **live events** that can be retrieved and updated dynamically.

### **Entities**

#### **Live Event**

- `id` (UUID)
- `title` (string)
- `description` (string)
- `start_time` (timestamp)
- `end_time` (timestamp)
- `rewards` (JSON string)

---

### **API Definitions**

#### **HTTP Endpoints (Public)**

| Method | Endpoint       | Description                  | Authentication               |
|--------|----------------|------------------------------|------------------------------|
| `GET`  | `/events`      | Retrieve all active events   | Requires API Key (HTTP role) |
| `GET`  | `/events/{id}` | Retrieve details of an event | Requires API Key (HTTP role) |

#### **gRPC Endpoints (Internal)**

| Method        | RPC Name      | Description                             | Authentication               |
|---------------|---------------|-----------------------------------------|------------------------------|
| `CreateEvent` | `CreateEvent` | Create a new live event                 | Requires API Key (gRPC role) |
| `UpdateEvent` | `UpdateEvent` | Update an existing event                | Requires API Key (gRPC role) |
| `DeleteEvent` | `DeleteEvent` | Delete an event                         | Requires API Key (gRPC role) |
| `ListEvents`  | `ListEvents`  | Retrieve all events (incl. past events) | Requires API Key (gRPC role) |

---

## **Authentication & Role-Based Access**

- API requests must include an **API Key** in the `Authorization` header.
- API keys have **roles**:
    - `"http_user"` → Can access **only HTTP endpoints**.
    - `"grpc_admin"` → Can access **only gRPC endpoints**.
- The system should **reject unauthorized access**.

---

## **Technical Specifications**

- **Single binary** running **both HTTP & gRPC servers on the same port**.
- Use **SQLite** for persistent storage.
- **Protobuf definitions** must be used for the gRPC API.
- The HTTP API should be served using **a RESTful framework**.
- The gRPC API should use **protobuf-generated Go code**.

---

## **Expected Deliverables**

1. **Golang Codebase**
    - Well-structured and maintainable code.
    - Proper separation of concerns.
    - Clean error handling and logging.

2. **README.md**
    - Instructions to build, run, and test the application.
    - API documentation (HTTP & gRPC).
    - Details on authentication and API Key roles.

3. **Dockerfile**
    - A Dockerfile to containerize the application.

4. **Unit Tests**
    - Test authentication logic.
    - Test core business logic (event retrieval, creation, updates).

---

## **Evaluation Criteria**

- **Code Structure & Best Practices**: Readability, maintainability, idiomatic Go code.
- **Concurrency Handling**: Efficient management of incoming HTTP & gRPC requests.
- **API Design**: Logical structuring of HTTP & gRPC endpoints.
- **Security Considerations**: Correct implementation of **API Key authentication**.
- **Database Handling**: Efficient SQLite usage.

---

## **Submission**

Please submit your solution as a **public Git repository** (GitHub, GitLab, etc.), including:

- The complete source code.
- A README with instructions.
- Any additional comments or documentation.
