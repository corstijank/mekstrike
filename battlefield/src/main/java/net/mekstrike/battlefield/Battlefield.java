package net.mekstrike.battlefield;

import java.util.List;

import io.dapr.actors.ActorMethod;
import io.dapr.actors.ActorType;

@ActorType(name = "battlefield")
public interface Battlefield {
    @ActorMethod(name = "GetNumberOfCols")
    public int getNumberOfCols();

    @ActorMethod(name = "GetNumberOfRows")
    public int getNumberOfRows();

    @ActorMethod(name = "IsCellBlocked")
    public boolean isCellBlocked(Coordinates cellRef);
    
    @ActorMethod(name = "BlockCell")
    public void blockCell(Coordinates cellRef);

    @ActorMethod(name="GetBoardCells")
    public List<Cell> getBoardCells();

    @ActorMethod(name="GetMovementOptions")
    public List<Coordinates> getMovementOptions(UnitData unit);
}
