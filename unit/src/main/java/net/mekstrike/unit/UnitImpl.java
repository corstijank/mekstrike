package net.mekstrike.unit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.dapr.actors.ActorId;
import io.dapr.actors.client.ActorClient;
import io.dapr.actors.client.ActorProxyBuilder;
import io.dapr.actors.runtime.AbstractActor;
import io.dapr.actors.runtime.ActorRuntimeContext;
import reactor.core.publisher.Mono;

public class UnitImpl extends AbstractActor implements Unit {
    private static final Logger LOGGER = LoggerFactory.getLogger(UnitImpl.class);

    private String owner;

    private UnitStats stats;

    private UnitLocation location;

    private ActorProxyBuilder<Battlefield> battlefieldBuilder;

    public UnitImpl(ActorRuntimeContext<UnitImpl> runtimeContext, ActorId id) {
        super(runtimeContext, id);
        battlefieldBuilder = new ActorProxyBuilder<Battlefield>(Battlefield.class, new ActorClient());
    }

    @Override
    public void deploy(DeployData data) {
        String battlefieldID = data.getBattlefieldID();
        owner = data.getOwner();
        stats = data.getStats();
        LOGGER.info("Deploying unit " + data.getStats().getModel() + " for player " + data.getOwner());

        var battlefield = battlefieldBuilder.build(new ActorId(battlefieldID));
        CellRef cell = null;
        // Default heading = UP
        // TODO: Extract heading to something less magically numbered
        int heading = 0;
        if (data.getDeployLocation().equalsIgnoreCase("NE")) {
            // start 0,0
            int row = 0;
            int col = 0;
            while (cell == null) {
                LOGGER.info("Checking if cell is blocked - Col " + col + " row " + row);
                if (!battlefield.isCellBlocked(new CellRef(col, row))) {
                    cell = new CellRef(col, row);
                } else {
                    col++;
                }
            }
            // Set heading to DOWN
            heading = 3;
        } else if (data.getDeployLocation().equalsIgnoreCase("SW")) {
            // Start bottomr right and walk back
            int row = battlefield.getNumberOfRows() - 1;
            int col = battlefield.getNumberOfCols() - 1;
            while (cell == null) {
                LOGGER.info("Checking if cell is blocked - Col " + col + " row " + row);
                if (!battlefield.isCellBlocked(new CellRef(col, row))) {
                    cell = new CellRef(col, row);
                } else {
                    col--;
                }
            }
        }
        LOGGER.info("Blocking cell for unit deployment - Col " + cell.getCol() + " row " + cell.getRow());
        battlefield.blockCell(cell);

        location = new UnitLocation(battlefieldID, cell, heading);

        LOGGER.info("Deployed unit " + data.getStats().getModel() + " for player " + data.getOwner() + " to row "
                + cell.getRow() + " col " + cell.getCol() + " of battlefield " + data.getBattlefieldID());

        var save = getActorStateManager().set("owner", owner)
                .and(getActorStateManager().set("stats", stats))
                .and(getActorStateManager().set("location", location));// Finish saving
        save.block();
    }

    @Override
    public UnitLocation getLocation() {
        return location;
    }

    /**
     * Callback function invoked after an Actor has been activated.
     *
     * @return Asynchronous void response.
     */
    protected Mono<Void> onActivate() {
        // TODO: I really should do async stuff better.
        if (getActorStateManager().contains("stats").block()) {
            owner = getActorStateManager().get("owner", String.class).block();
            stats = getActorStateManager().get("stats", UnitStats.class).block();
            location = getActorStateManager().get("location", UnitLocation.class).block();
        }
        return Mono.empty();
    }

    /**
     * Callback function invoked after an Actor has been deactivated.
     *
     * @return Asynchronous void response.
     */
    protected Mono<Void> onDeactivate() {
        return getActorStateManager().set("owner", owner)
                .and(getActorStateManager().set("stats", stats))
                .and(getActorStateManager().set("location", location));
    }

}
