using System;
using HexCore;
namespace Battlefield.Interfaces
{
    [Serializable]
    public class Hex
    {
        public int Col { get; set; }
        public int Row { get; set; }
        public int TerrainTypeID { get; set; }

        public static Hex FromCellState(CellState cs)
        {
            var c2d = cs.Coordinate3.To2D(OffsetTypes.OddRowsRight);
            return new Hex
            {
                Col = c2d.X,
                Row = c2d.Y,
                TerrainTypeID = cs.TerrainType.Id,
            };
        }
    }

}