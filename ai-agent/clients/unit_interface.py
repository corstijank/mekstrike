from typing import Dict, Any

class UnitActorInterface:
    async def GetData(self) -> Dict[str, Any]:
        pass
    
    async def SetActive(self, active: bool) -> None:
        pass