using HexCore;
using System.Collections.Generic;

namespace BattlefieldService.Board
{
    public static class Movement
    {
        public static MovementType Walking = new MovementType(1, "Walking");
        
        public static MovementTypes AllMovementTypes = new MovementTypes(
            new[] { Terrain.Ground },
            new Dictionary<MovementType, Dictionary<TerrainType, int>>
            {
                [Walking] = new Dictionary<TerrainType, int>
                {
                    [Terrain.Ground] = 1,
                }
            }
        );
    }
}