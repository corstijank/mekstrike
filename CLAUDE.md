# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Mekstrike is a distributed microservices-based tactical combat game built on Kubernetes. The system uses Dapr for service communication, gRPC/REST APIs, and supports a Svelte-based web UI. It implements a hexagonal battlefield with unit management across multiple programming languages.

## Architecture

### Service Components
- **gamemaster** (Go): Central game coordinator handling REST API endpoints for game state, player management, and turn coordination
- **armybuilder** (Java): gRPC service for army composition and force generation
- **battlefield** (Java): Actor-based service managing hex-based map logic and movement validation
- **unit** (Java): Actor-based service handling individual unit state, actions, and combat resolution
- **ai-agent** (Python): AI behavior service that processes AI unit turns using pub/sub events and Dapr service invocation
- **library** (Go): REST service providing unit data and statistics
- **importer** (Go): Data import service for unit definitions
- **mediaproxy** (Go): Asset serving proxy
- **ui** (Svelte): Web-based game interface

### Communication Patterns
- Services communicate via Dapr service invocation
- gRPC for high-performance service-to-service calls
- REST APIs for client-facing endpoints
- Protocol Buffers for data serialization
- Redis pub/sub for event-driven AI turn processing

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
./deploy.sh ai-agent
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

#### Python Services (ai-agent)
```bash
cd ai-agent
pip install -r requirements.txt  # Install dependencies
python app.py                    # Run service locally
```

## Key Development Concepts

### Dapr Integration
- All services are Dapr-enabled with sidecar architecture
- Service discovery uses Dapr service invocation
- State management through Dapr state stores (Redis)
- Configuration managed via Dapr configuration API
- Event-driven communication via Dapr pub/sub (Redis)

#### Messaging Patterns
**Service Invocation**: Direct API calls between services
- Go services call Java/Python services via Dapr HTTP/gRPC
- Example: `gamemaster` → `armybuilder` gRPC calls
- Example: `ai-agent` → `gamemaster` REST API calls

**Actor Model**: Stateful service instances with persistent state
- Unit actors maintain individual unit state and behavior
- Battlefield actors manage hex-grid state and validation
- Each actor has unique ID and persisted state in Redis

**Pub/Sub Events**: Asynchronous event-driven workflows  
- AI turn coordination via Redis pub/sub topics
- Future: Unit action events, combat resolution events
- Decoupled architecture allows easy service addition/modification

### Kubernetes Deployment
- Each service has corresponding Dockerfile and k8s yaml manifest in k8s/ directory
- Base platform components in k8s/mekstrike-base/ (namespace, config, stores)
- Platform dependencies in k8s/platform/ (monitoring, telemetry)
- Uses minikube for local development with image building

### Game Flow Architecture

#### Game Initialization
1. Player creates game via gamemaster REST API
2. Gamemaster coordinates with armybuilder to generate forces
3. Units are instantiated as independent actors via unit service
4. Battlefield manages hex-grid logic and movement validation
5. UI polls gamemaster for game state and renders board

#### Turn-Based Game Loop
The game follows a strict three-phase turn system:

**Phase Structure**: Movement → Combat → End → Next Unit
- Each unit in initiative order processes all three phases before advancing
- When all units complete a round, a new round begins with re-rolled initiative

**Turn Flow**:
1. **Movement Phase**: Unit moves on hex grid (AI or human player action)
2. **Combat Phase**: Unit declares and resolves attacks against targets
3. **End Phase**: Status effects, cleanup, unit removal if destroyed
4. **Advance**: Next unit activates or new round begins

#### AI vs Human Player Flow

**Human Players**:
- UI calls gamemaster APIs to get movement options, execute actions
- Player manually calls `POST /games/{id}/advanceTurn` after completing actions
- Gamemaster advances to next unit/phase/round

**AI Players (CPU)**:
- Gamemaster detects AI unit activation and publishes `ai-turn-started` event
- AI Agent receives event, analyzes game state via APIs, executes strategy
- AI Agent publishes phase-specific completion events automatically
- Gamemaster receives completion events and auto-advances turn

#### Event-Driven Messaging Architecture

**Pub/Sub Topics**:
- `ai-turn-started` - Gamemaster → AI Agent (when AI unit activated)
- `unit-movement-completed` - AI Agent → Gamemaster (movement phase done)
- `unit-attack-completed` - AI Agent → Gamemaster (combat phase done)  
- `unit-end-phase-completed` - AI Agent → Gamemaster (end phase done)

**Message Flow**:
```
┌─────────────┐    ai-turn-started    ┌──────────┐
│ Gamemaster  │──────────────────────→│ AI Agent │
│             │                       │          │
│   ┌─────────┼─ unit-*-completed ────┤          │
│   ↓         │←──────────────────────┤          │
│AdvanceTurn  │                       └──────────┘
└─────────────┘
```

This architecture allows seamless AI behavior while maintaining the same API contracts that human players use.

#### Key API Endpoints

**Gamemaster REST API** (`/games`):
- `POST /games` - Create new game with generated armies
- `GET /games/{id}` - Get game state (current phase, active unit, round)
- `GET /games/{id}/currentOpts` - Get movement options for active unit
- `GET /games/{id}/board` - Get battlefield with unit positions
- `GET /games/{id}/units/{uid}` - Get specific unit data
- `POST /games/{id}/advanceTurn` - Advance to next unit/phase (human players)

**Dapr Pub/Sub Events** (AI coordination):
- `GET /dapr/subscribe` - Dapr subscription configuration
- `POST /unit-movement-completed` - Handle AI movement completion
- `POST /unit-attack-completed` - Handle AI combat completion  
- `POST /unit-end-phase-completed` - Handle AI end phase completion

**Unit Actors** (Dapr actor invocation):
- `POST /v1.0/actors/unit/{id}/method/deploy` - Deploy unit to battlefield
- `POST /v1.0/actors/unit/{id}/method/setActive` - Activate/deactivate unit
- `GET /v1.0/actors/unit/{id}/method/getData` - Get unit state
- Future: `move`, `attack`, `declareTargets` methods

**Battlefield Actors** (Dapr actor invocation):
- `GET /v1.0/actors/battlefield/{id}/method/getBoardCells` - Get hex cells
- `POST /v1.0/actors/battlefield/{id}/method/getMovementOptions` - Valid moves
- `POST /v1.0/actors/battlefield/{id}/method/isCellBlocked` - Check occupation

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

## Important Corrections
- battlefield is Java, not dotnet, this was wrong in README