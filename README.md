# Technical Test Task – Backend Development

This small task is designed to assess your skills in developing web services. The task is intentionally simple but provides insight into your coding style, architectural decisions, and ability to implement requirements pragmatically.

## Task: Implementation of a Simple Calculator Web Service

Create a web service in Go (Golang) that supports basic arithmetic operations.

### Technical Specifications

The web service should provide the following operations:

| Operation | Input Parameters | Return Value |
|-----------|------------------|--------------|
| Addition | `float SummandOne`, `float SummandTwo` | `float Sum` |
| Subtraction | `float Minuend`, `float Subtrahend` | `float Difference` |
| Multiplication | `float FactorOne`, `float FactorTwo` | `float Product` |
| Division | `float Dividend`, `float Divisor` | `float Quotient` |
| Recent Results | optional: `int RecentN` (Default: 5, max: 20) | List of the last N calculations in any format (e.g., `["1 + 2 = 3", "3 * 5 = 15", ...]`) |

### Additional Requirements

- Recent calculations should be stored in-memory by default.
- Optionally, the web service should be startable with a persistence mode (e.g., via command-line parameter).
    - In this case, calculations should be stored in a local file and reloaded on startup.
    - The type of persistence (JSON file, SQLite, etc.) is up to you.

### Implementation Notes

- The technical framework (REST, gRPC, GraphQL, etc.) is freely selectable – choose what seems most appropriate to you.
- Ensure clean, understandable code structures.
- Logging, error handling, modularity, and API documentation (e.g., OpenAPI) are welcome but not mandatory.

### Goal

The goal is not a perfect or feature-rich service, but clearly structured, functional, and easily understandable code that you develop within a realistic timeframe.

### Submission

Please submit the source code as a repository (e.g., GitHub, GitLab, or ZIP) – ideally with a brief README on how to start and test the service.