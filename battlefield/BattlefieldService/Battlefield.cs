using Dapr.Actors;
using Dapr.Actors.Runtime;
using Battlefield.Interfaces;
using System;
using System.Threading.Tasks;
using System.Collections.Generic;
using BattlefieldService.Board;
using System.Text.Json;


namespace BattlefieldService
{
    [Actor(TypeName = "battlefield")]
    internal class Battlefield : Actor, IBattlefield
    {
        private Board.Board board;

        private int cols;
        private int rows;

        public Battlefield(ActorHost host) : base(host)
        {
            cols = 15;
            rows = 17;
            board = BoardGenerator.GenerateBlankMap(cols, rows);
        }

        public Task<List<Hex>> GetBoardCells()
        {
            var result = new List<Hex>();
            foreach (HexCore.CellState cs in board.MapGraph.CellsList)
            {
                result.Add(Hex.FromCellState(cs));
            }
            return Task.FromResult<List<Hex>>(result);
        }

        Task<int> IBattlefield.GetNumberOfCols()
        {
            return Task.FromResult<int>(cols);
        }

        Task<int> IBattlefield.GetNumberOfRows()
        {
            return Task.FromResult<int>(rows);
        }

        Task<bool> IBattlefield.IsCellBlocked(Position position)
        {
            return Task.FromResult<bool>(board.IsCellBlocked(position.Row, position.Col));
        }
        Task IBattlefield.BlockCell(Position position)
        {
            board.BlockCell(position.Row, position.Col);
            return Task.CompletedTask;
        }

        /// <summary>
        /// This method is called whenever an actor is activated.
        /// An actor is activated the first time any of its methods are invoked.
        /// </summary>
        protected override Task OnActivateAsync()
        {
            // Provides opportunity to perform some optional setup.
            Console.WriteLine($"Activating actor id: {this.Id}");
            return Task.CompletedTask;
        }

        /// <summary>
        /// This method is called whenever an actor is deactivated after a period of inactivity.
        /// </summary>
        protected override Task OnDeactivateAsync()
        {
            // Provides Oppertunity to perform optional cleanup.
            Console.WriteLine($"Deactivating actor id: {this.Id}");
            return Task.CompletedTask;
        }
    }
}