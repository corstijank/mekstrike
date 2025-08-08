import asyncio
import json
import logging
import nest_asyncio
from typing import Dict, Any
from fastapi import FastAPI, Request
from dapr.clients import DaprClient
from ai.strategy import AIStrategy
from clients.gamemaster_client import GamemasterClient
from clients.unit_client import UnitClient
from clients.battlefield_client import BattlefieldClient

# Apply nest_asyncio to allow nested event loops
nest_asyncio.apply()

# Setup logging
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

# FastAPI app for HTTP-based Dapr subscriptions
app = FastAPI(title="Mekstrike AI Agent")

# AI strategy engine
ai_strategy = AIStrategy()

# Global DaprClient instance
dapr_client = None

@app.on_event("startup")
async def startup_event():
    """Initialize resources on startup"""
    global dapr_client
    dapr_client = DaprClient(address='127.0.0.1:50001')

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "ai-agent"}

@app.get("/dapr/subscribe")
async def subscribe():
    """Return Dapr subscription configuration"""
    subscriptions = [
        {
            "pubsubname": "redis-pubsub",
            "topic": "ai-turn-started", 
            "route": "ai-turn-started"
        }
    ]
    return subscriptions

@app.post("/ai-turn-started")
async def ai_turn_handler(request: Request):
    """Handle AI turn events from gamemaster"""
    try:
        # Get the raw request body
        body = await request.body()
        data = json.loads(body)
       
        # Handle CloudEvent wrapper if present
        if 'data' in data:
            # The data field is a JSON string, parse it again
            if isinstance(data['data'], str):
                turn_data = json.loads(data['data'])
                logger.debug(f"Parsed data field as JSON: {turn_data}")
            else:
                turn_data = data['data']
                logger.debug(f"Using data field directly: {turn_data}")
        else:
            turn_data = data
            logger.debug(f"Using raw data: {turn_data}")
            
        logger.debug(f"Final turn_data type: {type(turn_data)}")
        logger.info(f"🎯 Processing AI turn event: {turn_data}")
        
        # Process the AI turn asynchronously
        await process_ai_turn(turn_data)
        
        logger.info(f"✅ Successfully processed AI turn for unit: {turn_data.get('unitId', 'unknown')}")
        return {"status": "success"}
        
    except Exception as e:
        logger.error(f"❌ Error processing AI turn: {e}", exc_info=True)
        return {"status": "error", "message": str(e)}

async def process_ai_turn(turn_data: Dict[str, Any]) -> None:
    """Process an AI turn based on the current game state"""
    global dapr_client
    
    game_id = turn_data.get('gameId')
    unit_id = turn_data.get('unitId')
    phase = turn_data.get('phase')
    round_num = turn_data.get('round')
    
    logger.info(f"Processing AI turn - Game: {game_id}, Unit: {unit_id}, Phase: {phase}, Round: {round_num}")
    
    try:
        # Use the global DaprClient instance
        if dapr_client is None:
            raise RuntimeError("DaprClient not initialized")
        
        # Initialize API clients
        gamemaster = GamemasterClient(dapr_client)
        unit_client = UnitClient(dapr_client)
        battlefield = BattlefieldClient(dapr_client)
        
        # Get current game state
        logger.info("Fetching game data...")
        game_data = await gamemaster.get_game(game_id)
        logger.info(f"Game data: {game_data}")
        
        logger.info("Fetching unit data...")
        unit_data = await unit_client.get_unit(unit_id)
        logger.info(f"Unit data: {unit_data}")
        
        logger.info("Fetching board data...")
        board_data = await gamemaster.get_board(game_id)
        logger.info(f"Board data: {board_data}")
        
        logger.info("Fetching move options...")
        move_options = await gamemaster.get_current_options(game_id)
        logger.info(f"Move options: {move_options}")
        
        # Execute AI strategy based on phase
        if phase == "Movement":
            await ai_strategy.execute_movement(
                dapr_client, game_data, unit_data, board_data, move_options, unit_id
            )
        elif phase == "Combat":
            await ai_strategy.execute_combat(
                dapr_client, game_data, unit_data, board_data, unit_id
            )
        elif phase == "End":
            await ai_strategy.execute_end_phase(
                dapr_client, game_data, unit_data, unit_id
            )
        
        logger.info(f"AI turn processing completed for {unit_id}")
            
    except Exception as e:
        logger.error(f"Error executing AI turn: {e}", exc_info=True)
        # TODO: Consider publishing error event

@app.on_event("shutdown")
async def shutdown_event():
    """Cleanup resources on shutdown"""
    global dapr_client
    if dapr_client is not None:
        await dapr_client.close()

if __name__ == "__main__":
    import uvicorn
    # Run FastAPI HTTP server for Dapr subscriptions
    uvicorn.run(app, host="0.0.0.0", port=50051)