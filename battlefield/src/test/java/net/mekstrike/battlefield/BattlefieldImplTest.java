package net.mekstrike.battlefield;

import io.dapr.actors.ActorId;
import io.dapr.actors.runtime.ActorTestUtils;
import net.mekstrike.domain.battlefield.Battlefield;
import net.mekstrike.domain.unit.Unit;
import org.hexworks.mixite.core.api.CubeCoordinate;
import org.hexworks.mixite.core.api.HexagonOrientation;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class BattlefieldImplTest {

    private BattlefieldImpl battlefield;

    @BeforeEach
    void setUp() {
        var context = ActorTestUtils.<BattlefieldImpl>createMockContext("battlefield");
        battlefield = new BattlefieldImpl(context, new ActorId("test-battlefield"));
    }

    @Test
    void getBoardCells_returnsAllCells() {
        var cells = battlefield.getBoardCells();

        // 17 rows x 15 cols = 255 hexagons
        assertEquals(255, cells.size());
    }

    @Test
    void getBoardCells_allCellsAreGrass() {
        var cells = battlefield.getBoardCells();

        assertTrue(cells.stream().allMatch(c -> "GRASS".equals(c.getType())));
    }

    @Test
    void getBoardCells_allCellsHavePositiveCoordinates() {
        var cells = battlefield.getBoardCells();

        assertTrue(cells.stream().allMatch(c ->
                c.getCoordinates().getX() >= 1 && c.getCoordinates().getY() >= 1));
    }

    @Test
    void getNumberOfCols_returns15() {
        assertEquals(15, battlefield.getNumberOfCols());
    }

    @Test
    void getNumberOfRows_returns17() {
        assertEquals(17, battlefield.getNumberOfRows());
    }

    @Test
    void isCellBlocked_returnsFalseForEmptyCell() {
        var coords = Battlefield.Coordinates.newBuilder().setX(3).setY(3).build();

        assertFalse(battlefield.isCellBlocked(coords));
    }

    @Test
    void blockCell_makesCellBlocked() {
        var coords = Battlefield.Coordinates.newBuilder().setX(5).setY(5).build();

        battlefield.blockCell(coords);

        assertTrue(battlefield.isCellBlocked(coords));
    }

    @Test
    void unblockCell_makesCellUnblocked() {
        var coords = Battlefield.Coordinates.newBuilder().setX(5).setY(5).build();
        battlefield.blockCell(coords);

        battlefield.unblockCell(coords);

        assertFalse(battlefield.isCellBlocked(coords));
    }

    @Test
    void getMovementOptions_returnsValidMoves() {
        var unit = createUnitAt(5, 5, "6\"");

        var options = battlefield.getMovementOptions(unit);

        assertFalse(options.isEmpty());
    }

    @Test
    void getMovementOptions_excludesBlockedCells() {
        var unit = createUnitAt(5, 5, "6\"");
        var allOptions = battlefield.getMovementOptions(unit);

        // Block one of the movement options
        var target = allOptions.getFirst();
        battlefield.blockCell(target);

        var optionsAfterBlock = battlefield.getMovementOptions(unit);

        assertTrue(optionsAfterBlock.size() < allOptions.size());
        assertFalse(optionsAfterBlock.contains(target));
    }

    @Test
    void getMovementOptions_moreMovementMeansMoreOptions() {
        var slowUnit = createUnitAt(8, 8, "4\"");
        var fastUnit = createUnitAt(8, 8, "12\"");

        var slowOptions = battlefield.getMovementOptions(slowUnit);
        var fastOptions = battlefield.getMovementOptions(fastUnit);

        assertTrue(fastOptions.size() > slowOptions.size());
    }

    @Test
    void convertCubeCoordinateToOffsetRow_flatTop() {
        var coord = CubeCoordinate.fromCoordinates(2, 3);

        int row = BattlefieldImpl.convertCubeCoordinateToOffsetRow(coord, HexagonOrientation.FLAT_TOP);

        assertEquals(1 + 3 + (2 - (2 & 1)) / 2, row);
    }

    @Test
    void convertCubeCoordinateToOffsetRow_pointyTop() {
        var coord = CubeCoordinate.fromCoordinates(2, 3);

        int row = BattlefieldImpl.convertCubeCoordinateToOffsetRow(coord, HexagonOrientation.POINTY_TOP);

        assertEquals(1 + 3, row);
    }

    @Test
    void convertCubeCoordinateToOffsetColumn_flatTop() {
        var coord = CubeCoordinate.fromCoordinates(2, 3);

        int col = BattlefieldImpl.convertCubeCoordinateToOffsetColumn(coord, HexagonOrientation.FLAT_TOP);

        assertEquals(1 + 2, col);
    }

    @Test
    void convertCubeCoordinateToOffsetColumn_pointyTop() {
        var coord = CubeCoordinate.fromCoordinates(2, 3);

        int col = BattlefieldImpl.convertCubeCoordinateToOffsetColumn(coord, HexagonOrientation.POINTY_TOP);

        assertEquals(1 + 2 + (3 - (3 & 1)) / 2, col);
    }

    private Unit.Data createUnitAt(int x, int y, String movement) {
        var position = Battlefield.Coordinates.newBuilder().setX(x).setY(y).build();
        var location = Unit.Location.newBuilder().setPosition(position).build();
        var stats = Unit.Stats.newBuilder().setMovement(movement).build();
        return Unit.Data.newBuilder()
                .setStats(stats)
                .setLocation(location)
                .build();
    }
}
