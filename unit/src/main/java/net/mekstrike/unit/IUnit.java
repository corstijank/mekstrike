package net.mekstrike.unit;

import io.dapr.actors.ActorMethod;
import io.dapr.actors.ActorType;
import net.mekstrike.domain.unit.Unit;

@ActorType(name = "unit")
public interface IUnit {
    @ActorMethod(name = "Deploy")
    void deploy(Unit.DeployRequest request);

    @ActorMethod(name = "GetData")
    Unit.Data getData();

    @ActorMethod(name = "SetActive")
    void setActive(boolean active);
}
