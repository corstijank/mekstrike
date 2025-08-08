# AI Agent Service

Python-based AI behavior service for Mekstrike tactical combat game.

## Overview

The AI Agent service processes AI unit turns using event-driven architecture:

1. **Event Subscription**: Listens for `ai-turn-started` events from gamemaster
2. **Decision Making**: Analyzes game state and executes AI strategies
3. **Action Execution**: Performs movement and combat actions via Dapr service calls
4. **Completion**: Publishes `ai-turn-completed` events

## Architecture

- **Python + FastAPI**: Web service framework
- **Dapr Pub/Sub**: Redis-based event messaging
- **Dapr Service Invocation**: API calls to gamemaster, unit, and battlefield services
- **Actor Communication**: Direct calls to unit and battlefield actors

## AI Strategy

Current implementation includes basic AI behaviors:

### Movement Phase
- Analyzes available movement options
- Simple positioning logic (move toward enemies, random valid moves)
- Respects battlefield constraints

### Combat Phase  
- Identifies targets within range
- Basic target selection (first available, priority to damaged units)
- Executes attacks via unit actors

### End Phase
- Cleanup and status effects processing
- Turn completion handling

## Development

```bash
# Install dependencies
pip install -r requirements.txt

# Run locally (requires Dapr sidecar)
python app.py

# Build and deploy
../deploy.sh ai-agent
```

## Event Flow

```
Gamemaster → [ai-turn-started] → AI Agent → API Calls → AI Agent → [unit-*-completed] → Gamemaster → Advance Turn
```

### Action-Specific Events
- `unit-movement-completed` - Published after movement phase processing
- `unit-attack-completed` - Published after combat phase processing  
- `unit-end-phase-completed` - Published after end phase processing

These events trigger automatic turn advancement in the gamemaster, creating a seamless AI turn flow.

## TODO

- Implement actual unit movement/attack methods in unit actors
- Add sophisticated AI strategies (line of sight, terrain analysis)
- Implement proper error handling and retry logic
- Add AI difficulty levels and personalities