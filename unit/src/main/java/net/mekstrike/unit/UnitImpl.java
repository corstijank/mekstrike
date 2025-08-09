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
import java.util.List;
import java.util.Map;
import java.util.HashMap;
import io.dapr.client.DaprClient;
import io.dapr.client.DaprClientBuilder;
import com.fasterxml.jackson.databind.ObjectMapper;

public class UnitImpl extends AbstractActor implements IUnit {
    private static final Logger logger = LoggerFactory.getLogger(UnitImpl.class);

    private String owner;

    private Boolean active;

    private Unit.Stats stats;

    private Unit.Location location;

    private ActorProxyBuilder<IBattlefield> battlefieldBuilder;
    private ActorClient actorClient;

    public UnitImpl(ActorRuntimeContext<UnitImpl> runtimeContext, ActorId id) {
        super(runtimeContext, id);
        actorClient = new ActorClient();
        battlefieldBuilder = new ActorProxyBuilder<IBattlefield>(IBattlefield.class, actorClient)
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

    @Override
    public String move(Object moveRequest) {
        // Handle case where moveRequest might be deserialized as a Map from JSON
        Battlefield.Coordinates actualTargetPosition;
        int newHeading = location.getHeading(); // Default to current heading
        
        if (moveRequest instanceof java.util.Map) {
            @SuppressWarnings("unchecked")
            java.util.Map<String, Object> requestMap = (java.util.Map<String, Object>) moveRequest;
            int x = ((Number) requestMap.get("x")).intValue();
            int y = ((Number) requestMap.get("y")).intValue();
            actualTargetPosition = Battlefield.Coordinates.newBuilder().setX(x).setY(y).build();
            
            // Extract heading if provided
            if (requestMap.containsKey("heading")) {
                newHeading = ((Number) requestMap.get("heading")).intValue();
            }
        } else if (moveRequest instanceof Battlefield.Coordinates) {
            actualTargetPosition = (Battlefield.Coordinates) moveRequest;
            // When using Coordinates directly, heading stays the same
        } else {
            throw new IllegalArgumentException("Invalid moveRequest type: " + moveRequest.getClass());
        }
        
        if (logger.isInfoEnabled()) {
            logger.info("Moving unit {} from x:{},y:{},heading:{} to x:{},y:{},heading:{}", 
                stats.getModel(), 
                location.getPosition().getX(), location.getPosition().getY(), location.getHeading(),
                actualTargetPosition.getX(), actualTargetPosition.getY(), newHeading);
        }
        
        // Validate movement is allowed
        var battlefield = battlefieldBuilder.build(new ActorId(location.getBattlefieldId()));
        var unitData = getData();
        List<?> allowedMoves = battlefield.getMovementOptions(unitData);
        
        boolean moveAllowed = allowedMoves.stream()
            .anyMatch(coord -> {
                // Handle case where coord might be deserialized as a Map from JSON
                if (coord instanceof java.util.Map) {
                    @SuppressWarnings("unchecked")
                    java.util.Map<String, Object> coordMap = (java.util.Map<String, Object>) coord;
                    int x = ((Number) coordMap.get("x")).intValue();
                    int y = ((Number) coordMap.get("y")).intValue();
                    return x == actualTargetPosition.getX() && y == actualTargetPosition.getY();
                } else if (coord instanceof Battlefield.Coordinates) {
                    Battlefield.Coordinates battleCoord = (Battlefield.Coordinates) coord;
                    return battleCoord.getX() == actualTargetPosition.getX() && battleCoord.getY() == actualTargetPosition.getY();
                } else {
                    logger.error("Unknown coordinate type in allowedMoves: {}", coord.getClass());
                    return false;
                }
            });
            
        if (!moveAllowed) {
            logger.error("Invalid movement attempt from x:{},y:{} to x:{},y:{} for unit {}", 
                location.getPosition().getX(), location.getPosition().getY(),
                actualTargetPosition.getX(), actualTargetPosition.getY(), stats.getModel());
            throw new IllegalArgumentException("Invalid movement coordinates");
        }
        
        // Store previous position for event publishing
        var previousPosition = location.getPosition();
        
        // Unblock previous cell
        battlefield.unblockCell(Battlefield.Coordinates.newBuilder()
            .setX(previousPosition.getX())
            .setY(previousPosition.getY())
            .build());
        
        // Block new cell  
        battlefield.blockCell(actualTargetPosition);
        
        // Update unit location
        location = Unit.Location.newBuilder()
            .setBattlefieldId(location.getBattlefieldId())
            .setPosition(actualTargetPosition)
            .setHeading(newHeading)
            .build();
            
        saveUnit().block();
        
        // Publish movement completed event
        publishMovementCompletedEvent(previousPosition, actualTargetPosition);
        
        if (logger.isInfoEnabled()) {
            logger.info("Successfully moved unit {} to x:{},y:{}", 
                stats.getModel(), actualTargetPosition.getX(), actualTargetPosition.getY());
        }
        
        return "{\"success\": true}";
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
        return saveUnit().doFinally(signal -> {
            // Close the ActorClient to prevent resource leaks
            if (actorClient != null) {
                try {
                    actorClient.close();
                } catch (Exception e) {
                    if (logger.isWarnEnabled()) {
                        logger.warn("Error closing ActorClient: {}", e.getMessage());
                    }
                }
            }
        });
    }

    private void publishMovementCompletedEvent(Battlefield.Coordinates previousPosition, Battlefield.Coordinates targetPosition) {
        DaprClient daprClient = null;
        try {
            daprClient = new DaprClientBuilder().build();
            // Create event data matching the AI agent's event format
            Map<String, Object> eventData = new HashMap<>();
            eventData.put("GameId", location.getBattlefieldId()); // GameId is same as BattlefieldId
            eventData.put("UnitId", getId().toString());
            eventData.put("Phase", "Movement");
            eventData.put("BattlefieldId", location.getBattlefieldId());
            
            Map<String, Integer> sourceLocation = new HashMap<>();
            sourceLocation.put("x", previousPosition.getX());
            sourceLocation.put("y", previousPosition.getY());
            eventData.put("SourceLocation", sourceLocation);
            
            Map<String, Integer> targetLocation = new HashMap<>();
            targetLocation.put("x", targetPosition.getX());
            targetLocation.put("y", targetPosition.getY());
            eventData.put("TargetLocation", targetLocation);
            
            Map<String, Object> unitInfo = new HashMap<>();
            unitInfo.put("Id", getId().toString());
            unitInfo.put("Model", stats.getModel());
            unitInfo.put("Owner", owner);
            
            Map<String, Integer> position = new HashMap<>();
            position.put("x", targetPosition.getX());
            position.put("y", targetPosition.getY());
            unitInfo.put("Position", position);
            unitInfo.put("Status", new HashMap<>());
            eventData.put("Unit", unitInfo);
            
            // Convert to JSON
            ObjectMapper mapper = new ObjectMapper();
            String jsonData = mapper.writeValueAsString(eventData);
            
            daprClient.publishEvent("redis-pubsub", "unit-movement-completed", jsonData).block();
            
            if (logger.isInfoEnabled()) {
                logger.info("Published unit-movement-completed event for unit {}", stats.getModel());
            }
        } catch (Exception e) {
            logger.error("Error publishing movement completed event for unit {}: {}", stats.getModel(), e.getMessage());
        } finally {
            if (daprClient != null) {
                try {
                    daprClient.close();
                } catch (Exception e) {
                    logger.warn("Error closing DaprClient: {}", e.getMessage());
                }
            }
        }
    }

    private Mono<Void> saveUnit() {
        return getActorStateManager().set("owner", owner)
                .and(getActorStateManager().set("stats", stats))
                .and(getActorStateManager().set("location", location))
                .and(getActorStateManager().set("active", active));
    }

}
