package net.mekstrike.battlefield;

import java.util.List;

import io.dapr.actors.ActorMethod;
import io.dapr.actors.ActorType;
import net.mekstrike.domain.battlefield.Battlefield;
import net.mekstrike.domain.unit.Unit;

@ActorType(name = "battlefield")
public interface IBattlefield {
    @ActorMethod(name = "GetNumberOfCols")
    int getNumberOfCols();

    @ActorMethod(name = "GetNumberOfRows")
    int getNumberOfRows();

    @ActorMethod(name = "IsCellBlocked")
    boolean isCellBlocked(Battlefield.Coordinates cellRef);
    
    @ActorMethod(name = "BlockCell")
    void blockCell(Battlefield.Coordinates cellRef);
    
    @ActorMethod(name = "UnblockCell")
    void unblockCell(Battlefield.Coordinates cellRef);

    @ActorMethod(name="GetBoardCells")
    List<Battlefield.Cell> getBoardCells();

    @ActorMethod(name="GetMovementOptions")
    List<Battlefield.Coordinates> getMovementOptions(Unit.Data unit);
}
