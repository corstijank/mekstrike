package net.mekstrike.unit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.dapr.actors.ActorId;
import io.dapr.actors.client.ActorClient;
import io.dapr.actors.client.ActorProxyBuilder;
import io.dapr.actors.runtime.AbstractActor;
import io.dapr.actors.runtime.ActorRuntimeContext;
import net.mekstrike.domain.battlefield.Battlefield;
import net.mekstrike.domain.unit.Unit;
import net.mekstrike.serialization.MekstrikeSerializer;
import net.mekstrike.unit.external.IBattlefield;
import reactor.core.publisher.Mono;

public class UnitImpl extends AbstractActor implements IUnit {
    private static final Logger logger = LoggerFactory.getLogger(UnitImpl.class);

    private String owner;

    private Boolean active;

    private Unit.Stats stats;

    private Unit.Location location;

    private ActorProxyBuilder<IBattlefield> battlefieldBuilder;

    public UnitImpl(ActorRuntimeContext<UnitImpl> runtimeContext, ActorId id) {
        super(runtimeContext, id);
        battlefieldBuilder = new ActorProxyBuilder<IBattlefield>(IBattlefield.class, new ActorClient())
                .withObjectSerializer(new MekstrikeSerializer());
    }

    @Override
    public void deploy(Unit.DeployRequest data) {
        String battlefieldID = data.getBattlefieldId();
        owner = data.getOwner();
        stats = data.getStats();
        active = false;

        if (logger.isInfoEnabled()) {
            logger.info("Deploying unit {} for player {}", data.getStats().getModel(), data.getOwner());
        }

        var battlefield = battlefieldBuilder.build(new ActorId(battlefieldID));
        Battlefield.Coordinates coordinates = null;

        // Default heading = UP
        // TODO: Extract heading to something less magically numbered
        int heading = 0;
        if ("NE".equalsIgnoreCase(data.getCorner())) {
            // start 0,0
            int y = 1;
            int x = 1;
            while (coordinates == null) {
                var c = Battlefield.Coordinates.newBuilder().setX(x).setY(y).build();
                logger.info("Checking if cell is blocked - x:{} y{} || {}", x, y, coordinates);

                if (!battlefield.isCellBlocked(c)) {
                    coordinates = c;
                } else {
                    x++;
                }
            }
            // Set heading to DOWN
            heading = 3;
        } else if ("SW".equalsIgnoreCase(data.getCorner())) {
            // Start bottomr right and walk back
            int y = battlefield.getNumberOfRows();
            int x = battlefield.getNumberOfCols();
            while (coordinates == null) {
                var c = Battlefield.Coordinates.newBuilder().setX(x).setY(y).build();
                logger.info("Checking if cell is blocked - x:{} y{} || {}", x, y, coordinates);

                if (!battlefield.isCellBlocked(c)) {
                    coordinates = c;
                } else {
                    x--;
                }
            }
        } else {
            throw new IllegalStateException(
                    "Deploy location for data " + data.getStats().getModel() + " should be NE or SW");
        }
        battlefield.blockCell(coordinates);

        location = Unit.Location.newBuilder().setBattlefieldId(battlefieldID).setPosition(coordinates)
                .setHeading(heading).build();

        if (logger.isInfoEnabled()) {
            logger.info("Deployed unit {} for player {} to x:{},y:{} of battlefield {}", data.getStats().getModel(),
                    owner,
                    coordinates.getX(), coordinates.getY(), battlefieldID);
        }
        saveUnit().block();
    }

    @Override
    public Unit.Data getData() {
        return Unit.Data.newBuilder().setLocation(location).setStats(stats).setOwner(owner).setActive(active).build();
    }

    @Override
    public void setActive(boolean active) {
        if (logger.isInfoEnabled()) {
            logger.info(" Activating unit {} for player {}", stats.getModel(), owner);
        }
        this.active = active;
        saveUnit().block();
    }

    /**
     * Callback function invoked after an Actor has been activated.
     *
     * @return Asynchronous void response.
     */
    protected Mono<Void> onActivate() {
        if (getActorStateManager().contains("stats").block()) {
            active = getActorStateManager().get("active",boolean.class).block();
            owner = getActorStateManager().get("owner", String.class).block();
            stats = getActorStateManager().get("stats", Unit.Stats.class).block();
            location = getActorStateManager().get("location", Unit.Location.class).block();
        }
        return Mono.empty();
    }

    /**
     * Callback function invoked after an Actor has been deactivated.
     *
     * @return Asynchronous void response.
     */
    protected Mono<Void> onDeactivate() {
        return saveUnit();
    }

    private Mono<Void> saveUnit() {
        return getActorStateManager().set("owner", owner)
                .and(getActorStateManager().set("stats", stats))
                .and(getActorStateManager().set("location", location))
                .and(getActorStateManager().set("active", active));
    }

}
