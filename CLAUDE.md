# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Building and Running
- `make run-agent` - Run the client agent application
- `make run-server` - Run the server application  
- `go run cmd/agent/main.go` - Run agent directly
- `go run cmd/server/main.go` - Run server directly

### Testing
- `go test ./...` - Run all tests
- `make testcoverage` - Run tests with coverage report (filters out generated files)
- `go test -v ./internal/path/to/package` - Run tests for specific package

### Code Generation
- `make proto` - Generate gRPC code from .proto files
- Uses protoc to generate Go code for users, items, and crypto services

### Database
- Uses SQLC for type-safe SQL code generation
- Schema files in `internal/server/repositories/database/schema/`
- Queries in `internal/server/repositories/database/query/`
- Generated code in `internal/server/repositories/database/generated/`

## Architecture Overview

GophKeeper is a client-server password manager with the following key components:

### Client-Server Communication
- Uses gRPC for communication between agent (client) and server
- Proto definitions in `internal/protos/` for users, items, and crypto services
- JWT tokens for authentication
- RSA encryption for secure data transmission

### Agent (Client) Architecture
- **UI Layer** (`internal/agent/ui/`): Terminal-based interface using state machine pattern
- **Services Layer** (`internal/agent/services/`): Business logic for users, items, and crypto
- **Client Layer** (`internal/agent/client/`): gRPC client implementations with interceptors
- Main entry point: `cmd/agent/main.go`

### Server Architecture  
- **Controllers** (`internal/server/controllers/`): gRPC service implementations
- **Services** (`internal/server/services/`): Business logic layer (user, item, crypto services)
- **Repository** (`internal/server/repositories/`): Data access layer with PostgreSQL
- **Crypto** (`internal/server/crypto/`): RSA key management and encryption utilities
- Main entry point: `cmd/server/main.go`

### Data Models
- **Item Types**: Credentials, Text, Binary, Card (defined in `models/item.go`)
- **Encryption**: All sensitive data is encrypted before storage
- **Database**: PostgreSQL with JSONB metadata support

### State Management (Agent UI)
- Complex state machine in `internal/agent/ui/states.go` with 50+ states
- Handles authentication, item management, and UI flows
- State transitions for login/signup, CRUD operations, and error handling

### Security Features
- RSA public/private key encryption
- Salted password hashing
- JWT-based authentication
- Client-side encryption before transmission
- Secure key storage and management

## Key Patterns

### Error Handling
- Consistent error wrapping with context
- Custom error types in `internal/errs/`
- Graceful error recovery in UI states

### Testing Strategy  
- Comprehensive test coverage with mocks
- Integration tests for database layer
- UI state machine testing
- Coverage filtering excludes generated code and protobuf files

### Configuration
- Environment-based configuration in `config/`
- Separate agent and server configurations
- Support for .env files via godotenv

### Logging
- Structured logging with zap logger
- Centralized logger initialization in `internal/logger/`