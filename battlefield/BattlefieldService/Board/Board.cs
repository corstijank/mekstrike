using System;
using HexCore;
using System.Collections.Generic;
namespace BattlefieldService.Board
{

    [Serializable]
    public class Board
    {
        public Graph MapGraph { get; set; }

        public Boolean IsCellBlocked(int row, int col)
        {
            var c = new Coordinate2D(x: col, y: row, OffsetTypes.OddRowsRight).To3D();
            return MapGraph.IsCellBlocked(c);
        }

        public void BlockCell(int row, int col)
        {
            var c = new Coordinate2D(x: col, y: row, OffsetTypes.OddRowsRight).To3D();
            MapGraph.BlockCells(new List<Coordinate3D>() { c });
        }
    }
}