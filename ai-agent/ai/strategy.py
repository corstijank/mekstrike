import logging
import math
import random
from typing import Dict, Any, List, Optional, Tuple
from dapr.clients import DaprClient
from clients.unit_client import UnitClient

logger = logging.getLogger(__name__)

class AIStrategy:
    """AI decision making engine for tactical combat"""
    
    def __init__(self):
        self.unit_client = None
    
    async def execute_movement(self, dapr_client: DaprClient, game_data: Dict[str, Any], 
                             unit_data: Dict[str, Any], board_data: Dict[str, Any], 
                             move_options: Dict[str, Any], unit_id: str) -> None:
        """Execute AI movement phase"""
        logger.info(f"AI executing movement for unit {unit_data.get('stats', {}).get('model', 'Unknown')} (ID: {unit_id})")
        
        try:
            self.unit_client = UnitClient(dapr_client)
            
            # Get available movement coordinates
            allowed_coords = move_options.get('AllowedCoordinates', [])
            if not allowed_coords:
                logger.info("No movement options available, skipping movement")
                await self._publish_movement_completed(dapr_client, game_data, unit_id)
                return
            
            # Find enemy units on the board
            enemy_units = self._find_enemy_units(game_data, unit_data, board_data)
            
            # Choose best movement position
            best_position = self._choose_movement_position(
                unit_data, allowed_coords, enemy_units, board_data
            )
            
            if best_position:
                # TODO: Execute movement once unit actor supports it
                logger.info(f"AI would move to position: {best_position}")
                # await self.unit_client.move_unit(unit_id, best_position)
            else:
                logger.info("No beneficial movement found, staying in place")
            
            # Publish movement completion event
            await self._publish_movement_completed(dapr_client, game_data, unit_id)
                
        except Exception as e:
            logger.error(f"Error executing AI movement: {e}")
            # Still publish completion to avoid hanging the game
            await self._publish_movement_completed(dapr_client, game_data, unit_id)
    
    async def execute_combat(self, dapr_client: DaprClient, game_data: Dict[str, Any],
                           unit_data: Dict[str, Any], board_data: Dict[str, Any], unit_id: str) -> None:
        """Execute AI combat phase"""
        logger.info(f"AI executing combat for unit {unit_data.get('stats', {}).get('model', 'Unknown')} (ID: {unit_id})")
        
        try:
            self.unit_client = UnitClient(dapr_client)
            
            # Find targets in range
            targets = self._find_targets_in_range(game_data, unit_data, board_data)
            
            if targets:
                # Choose best target
                best_target = self._choose_attack_target(unit_data, targets)
                
                if best_target:
                    # TODO: Execute attack once unit actor supports it
                    logger.info(f"AI would attack target: {best_target.get('id', 'Unknown')}")
                    # await self.unit_client.attack_unit(unit_id, best_target['id'])
                else:
                    logger.info("No suitable attack target found")
            else:
                logger.info("No targets in range for combat")
            
            # Publish attack completion event
            await self._publish_attack_completed(dapr_client, game_data, unit_id)
                
        except Exception as e:
            logger.error(f"Error executing AI combat: {e}")
            # Still publish completion to avoid hanging the game
            await self._publish_attack_completed(dapr_client, game_data, unit_id)
    
    async def execute_end_phase(self, dapr_client: DaprClient, game_data: Dict[str, Any],
                              unit_data: Dict[str, Any], unit_id: str) -> None:
        """Execute AI end phase"""
        logger.info(f"AI executing end phase for unit {unit_data.get('stats', {}).get('model', 'Unknown')} (ID: {unit_id})")
        
        try:
            # End phase typically involves cleanup, status effects, etc.
            # For now, just log that the turn is ending
            logger.info("AI end phase completed")
            
            # Publish end phase completion event
            await self._publish_end_phase_completed(dapr_client, game_data, unit_id)
            
        except Exception as e:
            logger.error(f"Error executing AI end phase: {e}")
            # Still publish completion to avoid hanging the game
            await self._publish_end_phase_completed(dapr_client, game_data, unit_id)
    
    def _find_enemy_units(self, game_data: Dict[str, Any], current_unit: Dict[str, Any], 
                         board_data: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Find enemy units on the battlefield"""
        enemy_units = []
        current_owner = current_unit.get('owner')
        
        # Get all cells with units
        cells = board_data.get('cells', [])
        for cell in cells:
            if cell.get('type') == 'unit':  # Assuming units are marked as type 'unit'
                # TODO: Need to get unit data from coordinates
                # This requires additional API calls to determine unit ownership
                pass
        
        # For now, return empty list - this needs proper implementation
        # once we can query units by position
        return enemy_units
    
    def _choose_movement_position(self, unit_data: Dict[str, Any], allowed_coords: List[Dict[str, int]],
                                enemy_units: List[Dict[str, Any]], board_data: Dict[str, Any]) -> Optional[Dict[str, int]]:
        """Choose the best movement position"""
        if not allowed_coords:
            return None
        
        current_pos = unit_data.get('location', {}).get('position', {})
        
        # Simple strategy: move towards center or random position for now
        # In a real implementation, this would consider:
        # - Distance to enemies
        # - Terrain advantages  
        # - Line of sight
        # - Defensive positions
        
        if len(allowed_coords) == 1:
            return allowed_coords[0]
        
        # For demo purposes, choose a random valid position
        # but prefer positions that are different from current position
        valid_moves = [
            coord for coord in allowed_coords 
            if coord.get('x') != current_pos.get('x') or coord.get('y') != current_pos.get('y')
        ]
        
        if valid_moves:
            return random.choice(valid_moves)
        else:
            return allowed_coords[0]
    
    def _find_targets_in_range(self, game_data: Dict[str, Any], unit_data: Dict[str, Any],
                              board_data: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Find enemy units within attack range"""
        targets = []
        
        # TODO: Implement range calculation based on unit stats
        # This requires:
        # - Unit weapon ranges (short/med/long)
        # - Distance calculation on hex grid
        # - Line of sight checks
        
        return targets
    
    def _choose_attack_target(self, unit_data: Dict[str, Any], targets: List[Dict[str, Any]]) -> Optional[Dict[str, Any]]:
        """Choose the best target to attack"""
        if not targets:
            return None
        
        # Simple strategy: attack the first available target
        # In a real implementation, this would consider:
        # - Target priority (damaged units, high-value targets)
        # - Damage potential
        # - Hit probability
        
        return targets[0]
    
    def _calculate_distance(self, pos1: Dict[str, int], pos2: Dict[str, int]) -> float:
        """Calculate distance between two hex positions"""
        x1, y1 = pos1.get('x', 0), pos1.get('y', 0)
        x2, y2 = pos2.get('x', 0), pos2.get('y', 0)
        
        # Simple Euclidean distance for now
        # Hex grids have more complex distance calculations
        return math.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2)
    
    async def _publish_movement_completed(self, dapr_client: DaprClient, game_data: Dict[str, Any], unit_id: str) -> None:
        """Publish unit movement completed event"""
        try:
            import json
            event_data = {
                "GameId": game_data.get('ID'),
                "UnitId": unit_id,
                "Phase": "Movement"
            }
            
            dapr_client.publish_event(
                pubsub_name='redis-pubsub',
                topic_name='unit-movement-completed',
                data=json.dumps(event_data)
            )
            
            logger.info(f"Published unit-movement-completed: {event_data}")
            
        except Exception as e:
            logger.error(f"Error publishing movement completed event: {e}")
    
    async def _publish_attack_completed(self, dapr_client: DaprClient, game_data: Dict[str, Any], unit_id: str) -> None:
        """Publish unit attack completed event"""
        try:
            import json
            event_data = {
                "GameId": game_data.get('ID'),
                "UnitId": unit_id,
                "Phase": "Combat"
            }
            
            dapr_client.publish_event(
                pubsub_name='redis-pubsub',
                topic_name='unit-attack-completed',
                data=json.dumps(event_data)
            )
            
            logger.info(f"Published unit-attack-completed: {event_data}")
            
        except Exception as e:
            logger.error(f"Error publishing attack completed event: {e}")
    
    async def _publish_end_phase_completed(self, dapr_client: DaprClient, game_data: Dict[str, Any], unit_id: str) -> None:
        """Publish unit end phase completed event"""
        try:
            import json
            event_data = {
                "GameId": game_data.get('ID'),
                "UnitId": unit_id,
                "Phase": "End"
            }
            
            dapr_client.publish_event(
                pubsub_name='redis-pubsub',
                topic_name='unit-end-phase-completed',
                data=json.dumps(event_data)
            )
            
            logger.info(f"Published unit-end-phase-completed: {event_data}")
            
        except Exception as e:
            logger.error(f"Error publishing end phase completed event: {e}")