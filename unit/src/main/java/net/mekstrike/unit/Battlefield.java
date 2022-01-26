package net.mekstrike.unit;

import io.dapr.actors.ActorMethod;
import io.dapr.actors.ActorType;

@ActorType(name = "battlefield")
public interface Battlefield {
    @ActorMethod(name = "GetNumberOfCols")
    public int getNumberOfCols();

    @ActorMethod(name = "GetNumberOfRows")
    public int getNumberOfRows();

    @ActorMethod(name = "IsCellBlocked")
    public boolean isCellBlocked(CellRef cellRef);
    
    @ActorMethod(name = "BlockCell")
    public void blockCell(CellRef cellRef);
}
