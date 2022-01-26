
using HexCore;

namespace BattlefieldService.Board{
    public static class BoardGenerator
    {
        public static Board GenerateBlankMap(int cols, int rows){
            Graph g = GraphFactory.CreateRectangularGraph(cols,rows,Movement.AllMovementTypes,Terrain.Ground,OffsetTypes.OddRowsRight);

            return new Board{MapGraph=g};
        }
    }
}