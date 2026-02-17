package net.mekstrike.unit;

import io.dapr.actors.ActorId;
import io.dapr.actors.client.ActorClient;
import io.dapr.actors.client.ActorProxyBuilder;
import io.dapr.actors.runtime.ActorTestUtils;
import net.mekstrike.domain.battlefield.Battlefield;
import net.mekstrike.domain.unit.Unit;
import net.mekstrike.unit.external.IBattlefield;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.lang.reflect.Field;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

class UnitImplTest {

    private UnitImpl unit;
    private IBattlefield mockBattlefield;

    @BeforeEach
    void setUp() throws Exception {
        // Mock ActorClient construction to prevent Dapr sidecar connection
        try (var ignored = Mockito.mockConstruction(ActorClient.class)) {
            var context = ActorTestUtils.<UnitImpl>createMockContext("unit");
            unit = new UnitImpl(context, new ActorId("test-unit"));
        }

        // Replace battlefieldBuilder with a mock that returns our mock battlefield
        mockBattlefield = mock(IBattlefield.class);
        var mockBuilder = mock(ActorProxyBuilder.class);
        when(mockBuilder.build(any(ActorId.class))).thenReturn(mockBattlefield);

        Field field = UnitImpl.class.getDeclaredField("battlefieldBuilder");
        field.setAccessible(true);
        field.set(unit, mockBuilder);
    }

    // --- deploy tests ---

    @Test
    void deploy_toNECorner_setsLocationAndHeading() {
        deployUnit("NE");

        var data = unit.getData();
        assertEquals(1, data.getLocation().getPosition().getX());
        assertEquals(1, data.getLocation().getPosition().getY());
        assertEquals(3, data.getLocation().getHeading()); // NE faces DOWN
    }

    @Test
    void deploy_toSWCorner_setsLocationAndHeading() {
        when(mockBattlefield.getNumberOfRows()).thenReturn(17);
        when(mockBattlefield.getNumberOfCols()).thenReturn(15);
        when(mockBattlefield.isCellBlocked(any())).thenReturn(false);

        deployUnit("SW");

        var data = unit.getData();
        assertEquals(15, data.getLocation().getPosition().getX());
        assertEquals(17, data.getLocation().getPosition().getY());
        assertEquals(0, data.getLocation().getHeading()); // SW faces UP
    }

    @Test
    void deploy_invalidCorner_throws() {
        assertThrows(IllegalStateException.class, () -> deployUnit("EAST"));
    }

    @Test
    void deploy_skipsBlockedCells() {
        var blockedCell = Battlefield.Coordinates.newBuilder().setX(1).setY(1).build();
        when(mockBattlefield.isCellBlocked(blockedCell)).thenReturn(true);
        // All other cells are unblocked (default returns false for unstubbed args)

        var stats = Unit.Stats.newBuilder().setModel("Atlas").setMovement("4\"").build();
        var request = Unit.DeployRequest.newBuilder()
                .setBattlefieldId("bf-1")
                .setOwner("player1")
                .setStats(stats)
                .setCorner("NE")
                .build();
        unit.deploy(request);

        var data = unit.getData();
        assertEquals(2, data.getLocation().getPosition().getX());
        assertEquals(1, data.getLocation().getPosition().getY());
    }

    @Test
    void deploy_blocksCellOnBattlefield() {
        deployUnit("NE");

        verify(mockBattlefield).blockCell(any(Battlefield.Coordinates.class));
    }

    // --- getData tests ---

    @Test
    void getData_afterDeploy_returnsCorrectData() {
        deployUnit("NE");

        var data = unit.getData();
        assertEquals("Atlas", data.getStats().getModel());
        assertEquals("player1", data.getOwner());
        assertFalse(data.getActive());
        assertNotNull(data.getLocation());
    }

    // --- setActive tests ---

    @Test
    void setActive_activatesUnit() {
        deployUnit("NE");

        unit.setActive(true);

        assertTrue(unit.getData().getActive());
    }

    @Test
    void setActive_deactivatesUnit() {
        deployUnit("NE");
        unit.setActive(true);

        unit.setActive(false);

        assertFalse(unit.getData().getActive());
    }

    // --- move tests ---

    @Test
    void move_withMapInput_updatesLocation() {
        deployUnit("NE");
        stubMovementOptions(3, 1);

        unit.move(Map.of("x", 3, "y", 1));

        var data = unit.getData();
        assertEquals(3, data.getLocation().getPosition().getX());
        assertEquals(1, data.getLocation().getPosition().getY());
    }

    @Test
    void move_withMapInputAndHeading_updatesHeading() {
        deployUnit("NE");
        stubMovementOptions(3, 1);

        unit.move(Map.of("x", 3, "y", 1, "heading", 5));

        assertEquals(5, unit.getData().getLocation().getHeading());
    }

    @Test
    void move_withCoordinatesInput_updatesLocation() {
        deployUnit("NE");
        stubMovementOptions(2, 2);

        var target = Battlefield.Coordinates.newBuilder().setX(2).setY(2).build();
        unit.move(target);

        var data = unit.getData();
        assertEquals(2, data.getLocation().getPosition().getX());
        assertEquals(2, data.getLocation().getPosition().getY());
    }

    @Test
    void move_invalidPosition_throws() {
        deployUnit("NE");
        // Return empty movement options - no valid moves
        when(mockBattlefield.getMovementOptions(any())).thenReturn(List.of());

        assertThrows(IllegalArgumentException.class, () ->
                unit.move(Map.of("x", 99, "y", 99)));
    }

    @Test
    void move_invalidType_throws() {
        deployUnit("NE");

        assertThrows(IllegalArgumentException.class, () ->
                unit.move("invalid"));
    }

    @Test
    void move_unblocksPreviousAndBlocksNewCell() {
        deployUnit("NE");
        // Reset invocation count from deploy
        clearInvocations(mockBattlefield);
        stubMovementOptions(3, 1);

        unit.move(Map.of("x", 3, "y", 1));

        // Should unblock previous position (1,1)
        verify(mockBattlefield).unblockCell(
                Battlefield.Coordinates.newBuilder().setX(1).setY(1).build());
        // Should block new position (3,1)
        verify(mockBattlefield).blockCell(
                Battlefield.Coordinates.newBuilder().setX(3).setY(1).build());
    }

    @Test
    void move_returnsSuccessJson() {
        deployUnit("NE");
        stubMovementOptions(3, 1);

        var result = unit.move(Map.of("x", 3, "y", 1));

        assertEquals("{\"success\": true}", result);
    }

    // --- helpers ---

    private void deployUnit(String corner) {
        when(mockBattlefield.isCellBlocked(any())).thenReturn(false);
        when(mockBattlefield.getNumberOfRows()).thenReturn(17);
        when(mockBattlefield.getNumberOfCols()).thenReturn(15);

        var stats = Unit.Stats.newBuilder()
                .setModel("Atlas")
                .setMovement("4\"")
                .build();
        var request = Unit.DeployRequest.newBuilder()
                .setBattlefieldId("bf-1")
                .setOwner("player1")
                .setStats(stats)
                .setCorner(corner)
                .build();

        unit.deploy(request);
    }

    private void stubMovementOptions(int targetX, int targetY) {
        var target = Battlefield.Coordinates.newBuilder()
                .setX(targetX).setY(targetY).build();
        when(mockBattlefield.getMovementOptions(any())).thenReturn(List.of(target));
    }
}
