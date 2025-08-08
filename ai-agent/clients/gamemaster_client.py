import json
import logging
from typing import Dict, Any, List
from dapr.clients import DaprClient

logger = logging.getLogger(__name__)

class GamemasterClient:
    """Client for interacting with gamemaster service via Dapr"""
    
    def __init__(self, dapr_client: DaprClient):
        self.dapr_client = dapr_client
        self.service_name = "gamemaster"
    
    async def get_game(self, game_id: str) -> Dict[str, Any]:
        """Get game state"""
        try:
            response = await self.dapr_client.invoke_method_async(
                app_id=self.service_name,
                method_name=f"games/{game_id}",
                data=b'',
                http_verb="GET"
            )
            
            if isinstance(response.data, str):
                return json.loads(response.data)
            elif isinstance(response.data, bytes):
                return json.loads(response.data.decode('utf-8'))
            else:
                # Already parsed
                return response.data
        except Exception as e:
            logger.error(f"Error getting game {game_id}: {e}")
            raise
    
    async def get_current_options(self, game_id: str) -> Dict[str, Any]:
        """Get current movement options for active unit"""
        try:
            response = await self.dapr_client.invoke_method_async(
                app_id=self.service_name,
                method_name=f"games/{game_id}/currentOpts",
                data=b'',
                http_verb="GET"
            )
            return json.loads(response.data)
        except Exception as e:
            logger.error(f"Error getting current options for game {game_id}: {e}")
            raise
    
    async def get_board(self, game_id: str) -> Dict[str, Any]:
        """Get board state with cells, rows, cols"""
        try:
            response = await self.dapr_client.invoke_method_async(
                app_id=self.service_name,
                method_name=f"games/{game_id}/board",
                data=b'',
                http_verb="GET"
            )
            return json.loads(response.data)
        except Exception as e:
            logger.error(f"Error getting board for game {game_id}: {e}")
            raise
    
    async def get_unit(self, game_id: str, unit_id: str) -> Dict[str, Any]:
        """Get specific unit data"""
        try:
            response = await self.dapr_client.invoke_method_async(
                app_id=self.service_name,
                method_name=f"games/{game_id}/units/{unit_id}",
                data=b'',
                http_verb="GET"
            )
            return json.loads(response.data)
        except Exception as e:
            logger.error(f"Error getting unit {unit_id} in game {game_id}: {e}")
            raise