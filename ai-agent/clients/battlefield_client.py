import json
import logging
from typing import Dict, Any, List
from dapr.clients import DaprClient

logger = logging.getLogger(__name__)

class BattlefieldClient:
    """Client for interacting with battlefield actors via Dapr"""
    
    def __init__(self, dapr_client: DaprClient):
        self.dapr_client = dapr_client
        self.actor_type = "battlefield"
    
    async def get_movement_options(self, battlefield_id: str, unit_data: Dict[str, Any]) -> List[Dict[str, int]]:
        """Get available movement coordinates for unit"""
        try:
            response = await self.dapr_client.invoke_method_async(
                app_id="battlefield",
                method_name=f"actors/{self.actor_type}/{battlefield_id}/method/GetMovementOptions",
                data=json.dumps(unit_data).encode('utf-8'),
                http_verb="POST"
            )
            return json.loads(response.data) if response.data else []
        except Exception as e:
            logger.error(f"Error getting movement options for battlefield {battlefield_id}: {e}")
            raise
    
    async def get_board_cells(self, battlefield_id: str) -> List[Dict[str, Any]]:
        """Get all board cells"""
        try:
            response = await self.dapr_client.invoke_method_async(
                app_id="battlefield",
                method_name=f"actors/{self.actor_type}/{battlefield_id}/method/GetBoardCells",
                data=b'',
                http_verb="POST"
            )
            return json.loads(response.data) if response.data else []
        except Exception as e:
            logger.error(f"Error getting board cells for battlefield {battlefield_id}: {e}")
            raise
    
    async def is_cell_blocked(self, battlefield_id: str, coordinates: Dict[str, int]) -> bool:
        """Check if cell is blocked"""
        try:
            response = await self.dapr_client.invoke_method_async(
                app_id="battlefield",
                method_name=f"actors/{self.actor_type}/{battlefield_id}/method/IsCellBlocked",
                data=json.dumps(coordinates).encode('utf-8'),
                http_verb="POST"
            )
            return json.loads(response.data) if response.data else False
        except Exception as e:
            logger.error(f"Error checking if cell blocked: {e}")
            raise
    
    async def get_number_of_rows(self, battlefield_id: str) -> int:
        """Get number of rows in battlefield"""
        try:
            response = await self.dapr_client.invoke_method_async(
                app_id="battlefield",
                method_name=f"actors/{self.actor_type}/{battlefield_id}/method/GetNumberOfRows",
                data=b'',
                http_verb="POST"
            )
            return json.loads(response.data) if response.data else 0
        except Exception as e:
            logger.error(f"Error getting number of rows: {e}")
            raise
    
    async def get_number_of_cols(self, battlefield_id: str) -> int:
        """Get number of columns in battlefield"""
        try:
            response = await self.dapr_client.invoke_method_async(
                app_id="battlefield",
                method_name=f"actors/{self.actor_type}/{battlefield_id}/method/GetNumberOfCols",
                data=b'',
                http_verb="POST"
            )
            return json.loads(response.data) if response.data else 0
        except Exception as e:
            logger.error(f"Error getting number of cols: {e}")
            raise