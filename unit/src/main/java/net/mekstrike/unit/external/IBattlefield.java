package net.mekstrike.unit.external;

import io.dapr.actors.ActorMethod;
import io.dapr.actors.ActorType;
import net.mekstrike.domain.battlefield.Battlefield;

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
}
