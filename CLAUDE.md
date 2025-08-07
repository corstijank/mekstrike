# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Mekstrike is a distributed microservices-based tactical combat game built on Kubernetes. The system uses Dapr for service communication, gRPC/REST APIs, and supports a Svelte-based web UI. It implements a hexagonal battlefield with unit management across multiple programming languages.

## Architecture

### Service Components
- **gamemaster** (Go): Central game coordinator handling REST API endpoints for game state, player management, and turn coordination
- **armybuilder** (Java): gRPC service for army composition and force generation
- **battlefield** (.NET): Actor-based service managing hex-based map logic and movement validation
- **unit** (Java): Actor-based service handling individual unit state, actions, and combat resolution
- **library** (Go): REST service providing unit data and statistics
- **importer** (Go): Data import service for unit definitions
- **mediaproxy** (Go): Asset serving proxy
- **ui** (Svelte): Web-based game interface

### Communication Patterns
- Services communicate via Dapr service invocation
- gRPC for high-performance service-to-service calls
- REST APIs for client-facing endpoints
- Protocol Buffers for data serialization

## Development Commands

### Platform Setup
```bash
# Initial kubernetes platform setup (Dapr, Jaeger, OpenTelemetry, Redis, Traefik)
./setup.sh

# Build all services and deploy to kubernetes
./deploy.sh

# Build and deploy specific service
./deploy.sh ui
./deploy.sh gamemaster
./deploy.sh armybuilder
# etc.
```

### Protocol Buffer Generation
```bash
# Generate Go protobuf files (run from project root)
protoc --proto_path="." --go_out=. --go_opt=paths=source_relative domain/unit/unit.proto
protoc --proto_path="." --go_out=. --go_opt=paths=source_relative domain/battlefield/battlefield.proto
protoc --proto_path="." --go_out=gamemaster/clients/armybuilder armybuilder/src/main/proto/armybuilder.proto
```

### Service Development

#### UI (Svelte)
```bash
cd ui
npm run dev        # Development server
npm run build      # Production build
npm run lint       # Lint code
npm run format     # Format code
npm run check      # Type checking
```

#### Java Services (armybuilder, unit)
```bash
cd armybuilder  # or unit
mvn exec:java      # Run service locally
mvn clean compile # Build
```

#### Go Services (gamemaster, library, importer, mediaproxy)
```bash
cd gamemaster  # or other go service
go run .           # Run service locally
go build           # Build binary
```

#### .NET Services (battlefield)
```bash
cd battlefield
dotnet run         # Run service locally
dotnet build       # Build
```

## Key Development Concepts

### Dapr Integration
- All services are Dapr-enabled with sidecar architecture
- Service discovery uses Dapr service invocation
- State management through Dapr state stores (Redis)
- Configuration managed via Dapr configuration API

### Kubernetes Deployment
- Each service has corresponding Dockerfile and k8s yaml manifest in k8s/ directory
- Base platform components in k8s/mekstrike-base/ (namespace, config, stores)
- Platform dependencies in k8s/platform/ (monitoring, telemetry)
- Uses minikube for local development with image building

### Game Flow Architecture
1. Player creates game via gamemaster REST API
2. Gamemaster coordinates with armybuilder to generate forces
3. Units are instantiated as independent actors
4. Battlefield manages hex-grid logic and movement validation
5. UI polls gamemaster for game state and renders board

### Observability
- OpenTelemetry instrumentation across services
- Jaeger distributed tracing
- Structured logging with service correlation

## File Structure Patterns
- Each service has independent Dockerfile for containerization
- Java services use Maven with exec plugin for local development
- Go services use standard module structure with domain-driven organization
- Protobuf definitions in domain/ and service-specific directories

## Development Notes
- When working on services, use `./deploy.sh` script for all build/deploy tasks, as local instructions using `mvnw` etc. may not work without Dapr

## Tool Insights
- technically; you can use the tools for compile, automated testing, but...like...don't expect to any integration beyond the unit working on