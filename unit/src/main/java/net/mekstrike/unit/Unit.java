package net.mekstrike.unit;

import io.dapr.actors.ActorMethod;
import io.dapr.actors.ActorType;

@ActorType(name = "unit")
public interface Unit {
    @ActorMethod(name = "Deploy")
    public void deploy(DeployData data);

    @ActorMethod(name = "GetLocation")
    public UnitLocation getLocation();

    @ActorMethod(name = "GetData")
    public UnitData getData();
}
