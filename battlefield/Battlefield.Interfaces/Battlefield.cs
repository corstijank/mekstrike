using Dapr.Actors;
using System.Threading.Tasks;
using System.Collections.Generic;

namespace Battlefield.Interfaces
{
    public interface IBattlefield : IActor
    {
        Task<List<Hex>> GetBoardCells();

        Task<List<Hex>> GetMovementOptions(int startX, int startY, string movement);

        Task<int>GetNumberOfCols();

        Task<int>GetNumberOfRows();

        Task<bool> IsCellBlocked(Position position);
        Task BlockCell(Position position);
    }
}