# Calculator Web Service

A REST API service in Go that provides basic arithmetic operations with calculation history storage.

## Features

- **Arithmetic Operations**: Addition, subtraction, multiplication, division
- **Storage Options**: In-memory or file-based persistence
- **Recent Calculations**: Retrieve last N calculations (default: 5, max: 20)
- **Metrics**: Prometheus metrics endpoint
- **Logging**: Structured logging with configurable levels

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/calculate/addition` | Add two numbers |
| POST | `/calculate/subtraction` | Subtract two numbers |
| POST | `/calculate/multiplication` | Multiply two numbers |
| POST | `/calculate/division` | Divide two numbers |
| GET | `/calculate/recent` | Get recent calculations |
| GET | `/metrics` | Prometheus metrics |

### Request Format
```json
{
  "operand1": 10.5,
  "operand2": 2.5
}
```

### Response Format
```json
{
  "result": 13.0,
  "operation": "addition",
  "expression": "10.5 + 2.5 = 13.0"
}
```

## Local Development

### Prerequisites
- Go 1.25+
- Make (optional)

### Quick Start
```bash
# Clone and navigate to project
git clone https://github.com/Chichichkin/CalculatorWebService.git
cd CalculatorWebService

# Run with memory storage (default)
go run ./cmd/main.go

# Run with file storage
export CALCULATOR_STORAGE_TYPE=file
export CALCULATOR_STORAGE_PATH=%prefered path% # Default ./storage.txt
touch storage.txt  # Create file if it doesn't exist. It will be created automatically on shutdown.

go run ./cmd/main.go
```

### Environment Variables
Here is the example of environment variables you can set:
```bash
CALCULATOR_PORT=8080                    # Server port
CALCULATOR_STORAGE_TYPE=memory          # Storage type: memory|file
CALCULATOR_STORAGE_PATH=./storage.txt   # File path for file storage
LOG_LEVEL=info                          # Log level: debug|info|warn|error
LOG_FORMAT=text                         # Log format: text|json
```
Rest can be found in `config/config.go` or `docker-compose.yml`
### Testing
```bash
# Test addition
curl -X POST http://localhost:8080/calculate/addition \
  -H "Content-Type: application/json" \
  -d '{"operand1": 5, "operand2": 3}'

# Get recent calculations
curl http://localhost:8080/calculate/recent

# Check metrics
curl http://localhost:8080/metrics
```

## Docker Deployment

### Using Docker Compose (Recommended)
```bash
# Start services
make compose-up

# Services available at:
# - Calculator (memory): http://localhost:8080
# - Calculator (file): http://localhost:8081
# - Prometheus: http://localhost:9090

# Run tests
make test-api
make test-memory
make test-file

# Stop services
make compose-down
```

### Manual Docker Build
```bash
# Build image
docker build -t calculator-service .

# Run with memory storage
docker run -p 8080:8080 calculator-service

# Run with file storage
docker run -p 8080:8080 \
  -e CALCULATOR_STORAGE_TYPE=file \
  -e CALCULATOR_STORAGE_PATH=/app/storage/calculations.txt \
  -v $(pwd)/storage:/app/storage \
  calculator-service
```

## Available Make Commands

```bash
make help              # Show all commands
make build             # Build Go binary
make run               # Run locally
make test              # Run Go tests
make docker-build      # Build Docker image
make compose-up        # Start with Docker Compose
make compose-down      # Stop services
make test-api          # Test API endpoints
make test-memory       # Test memory storage
make test-file         # Test file storage
make health            # Check service health
```
