import logging
from typing import Dict, Any
from dapr.actor import ActorProxy, ActorId
from .unit_interface import UnitActorInterface

logger = logging.getLogger(__name__)

class UnitClient:
    """Client for interacting with unit actors via Dapr"""
    
    def __init__(self, dapr_client=None):
        # We don't need dapr_client for ActorProxy approach
        self.actor_type = "unit"
    
    async def get_unit(self, unit_id: str) -> Dict[str, Any]:
        """Get unit data from unit actor"""
        try:
            # Create actor proxy
            proxy = ActorProxy.create(self.actor_type, ActorId(unit_id), UnitActorInterface)
            # Call the actor method
            data = await proxy.GetData()
            return data if isinstance(data, dict) else {}
        except Exception as e:
            logger.error(f"Error getting unit data for {unit_id}: {e}")
            raise
    
    async def set_active(self, unit_id: str, active: bool) -> None:
        """Set unit active state"""
        try:
            # Create actor proxy
            proxy = ActorProxy.create(self.actor_type, ActorId(unit_id), UnitActorInterface)
            # Call the actor method
            await proxy.SetActive(active)
            logger.info(f"Set unit {unit_id} active: {active}")
        except Exception as e:
            logger.error(f"Error setting unit {unit_id} active={active}: {e}")
            raise
    
    async def move_unit(self, unit_id: str, coordinates: Dict[str, int], heading: int = None) -> None:
        """Move unit to specified coordinates with optional heading"""
        try:
            # Create move request with coordinates and heading
            move_request = {
                "x": coordinates["x"],
                "y": coordinates["y"]
            }
            
            # Add heading if provided
            if heading is not None:
                move_request["heading"] = heading
                
            logger.info(f"Moving unit {unit_id} to {move_request}")
            
            # Create actor proxy
            proxy = ActorProxy.create(self.actor_type, ActorId(unit_id), UnitActorInterface)
            # Call the Move method with move request
            await proxy.Move(move_request)
            logger.info(f"Successfully moved unit {unit_id} to {move_request}")
        except Exception as e:
            logger.error(f"Error moving unit {unit_id}: {e}")
            raise
    
    async def attack_unit(self, unit_id: str, target_id: str) -> None:
        """Execute attack against target unit"""
        try:
            # TODO: Implement attack method once unit actor supports it
            logger.info(f"Unit {unit_id} attacking {target_id}")
            # This will need to be implemented in the unit actor
            pass
        except Exception as e:
            logger.error(f"Error unit {unit_id} attacking {target_id}: {e}")
            raise