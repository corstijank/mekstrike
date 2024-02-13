package net.mekstrike.battlefield;

import io.dapr.actors.ActorId;
import io.dapr.actors.runtime.AbstractActor;
import io.dapr.actors.runtime.ActorRuntimeContext;

import net.mekstrike.domain.battlefield.Battlefield;
import net.mekstrike.domain.unit.Unit;

import org.hexworks.mixite.core.api.CoordinateConverter;
import org.hexworks.mixite.core.api.CubeCoordinate;
import org.hexworks.mixite.core.api.Hexagon;
import org.hexworks.mixite.core.api.HexagonOrientation;
import org.hexworks.mixite.core.api.HexagonalGrid;
import org.hexworks.mixite.core.api.HexagonalGridBuilder;
import org.hexworks.mixite.core.api.HexagonalGridLayout;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import reactor.core.publisher.Mono;

import java.util.ArrayList;
import java.util.List;

public class BattlefieldImpl extends AbstractActor implements IBattlefield {
    private static final Logger LOGGER = LoggerFactory.getLogger(BattlefieldImpl.class);

    HexagonalGrid<BattlefieldHexData> grid;

    public BattlefieldImpl(ActorRuntimeContext<BattlefieldImpl> runtimeContext, ActorId id) {
        super(runtimeContext, id);
        HexagonalGridBuilder<BattlefieldHexData> builder = new HexagonalGridBuilder<BattlefieldHexData>()
                .setGridHeight(17)
                .setGridWidth(15)
                .setGridLayout(HexagonalGridLayout.RECTANGULAR)
                .setOrientation(HexagonOrientation.FLAT_TOP)
                .setRadius(30.0);
        grid = builder.build();
        saveBattlefield().block();
    }

    @Override
    public List<Battlefield.Cell> getBoardCells() {
        LOGGER.info("GetBoardCells called");

        var result = new ArrayList<Battlefield.Cell>();
        for (Hexagon<BattlefieldHexData> hexagon : grid.getHexagons()) {
            var c = hexagon.getCubeCoordinate();
            var x = convertCubeCoordinateToOffsetColumn(c, HexagonOrientation.FLAT_TOP);
            var y = convertCubeCoordinateToOffsetRow(c, HexagonOrientation.FLAT_TOP);
            var coordinates = Battlefield.Coordinates.newBuilder().setX(x).setY(y).build();
            var cell = Battlefield.Cell.newBuilder().setCoordinates(coordinates).setType("GRASS").build();
            result.add(cell);
        }
        return result;
    }

    @Override
    public int getNumberOfCols() {
        return grid.getGridData().getGridWidth();
    }

    @Override
    public int getNumberOfRows() {
        return grid.getGridData().getGridHeight();
    }

    @Override
    public boolean isCellBlocked(Battlefield.Coordinates cellRef) {
        LOGGER.info("isCellBlocked called");

        var h = getHex(cellRef.getX(),
                cellRef.getY());
        if (!h.getSatelliteData().isPresent()) {
            h.setSatelliteData(new BattlefieldHexData());
        }
        return h.getSatelliteData().get().isOccupied();
    }

    @Override
    public void blockCell(Battlefield.Coordinates cellRef) {
        LOGGER.info("BlockCell called");
        var h = getHex(cellRef.getX(),
                cellRef.getY());
        if (!h.getSatelliteData().isPresent()) {
            h.setSatelliteData(new BattlefieldHexData());
        }
        h.getSatelliteData().get().setOccupied(true);

        saveBattlefield().block();
    }

    @Override
    public List<Battlefield.Coordinates> getMovementOptions(Unit.Data unit) {
        LOGGER.info("GetMovementOptions called");
        var result = new ArrayList<Battlefield.Coordinates>();

        var calculator = new HexagonalGridBuilder<BattlefieldHexData>().buildCalculatorFor(grid);

        var h = getHex(unit.getLocation().getPosition().getX(),
                unit.getLocation().getPosition().getY());
        var movement = Integer.parseInt(unit.getStats().getMovement().split("\"")[0]) / 2;
        var options = calculator.calculateMovementRangeFrom(h, movement);
        for (Hexagon<BattlefieldHexData> hexagon : options) {
            if (hexagon.getSatelliteData().isPresent() && !hexagon.getSatelliteData().get().isOccupied()
                    || !hexagon.getSatelliteData().isPresent()) {
                var oc = hexagon.getCubeCoordinate();
                var ox = convertCubeCoordinateToOffsetColumn(oc, HexagonOrientation.FLAT_TOP);
                var oy = convertCubeCoordinateToOffsetRow(oc, HexagonOrientation.FLAT_TOP);
                result.add(Battlefield.Coordinates.newBuilder().setX(ox).setY(oy).build());
            }
        }

        return result;
    }

    public Mono<Void> saveBattlefield() {
        LOGGER.info("SaveBattlefield called");
        return getActorStateManager().set("grid", grid);
    }

    private Hexagon<BattlefieldHexData> getHex(int inx, int iny) {
        var x = CoordinateConverter.convertOffsetCoordinatesToCubeX(inx-1,
                iny-1, HexagonOrientation.FLAT_TOP);
        var z = CoordinateConverter.convertOffsetCoordinatesToCubeZ(inx-1,
                iny-1, HexagonOrientation.FLAT_TOP);
        CubeCoordinate c = CubeCoordinate.fromCoordinates(x, z);
        var maybeHex = grid.getByCubeCoordinate(c);
        if (maybeHex.isPresent()) {
            return maybeHex.get();
        }
        throw new IllegalStateException("No hex found for " + inx + "," + iny);
    }

    @SuppressWarnings("unchecked")
    @Override
    protected Mono<Void> onActivate() {
        if (Boolean.TRUE.equals(getActorStateManager().contains("grid").block())) {

            grid = getActorStateManager().get("grid", HexagonalGrid.class).block();
        }
        return Mono.empty();
    }

    @Override
    protected Mono<Void> onDeactivate() {
        return saveBattlefield();
    }

    /**
     * Calculates the offset row based on a CubeCoordinate.
     *
     * @param coordinate  a cube coordinate
     * @param orientation orientation
     *
     * @return offset row or y-value
     */
    public static int convertCubeCoordinateToOffsetRow(CubeCoordinate coordinate, HexagonOrientation orientation) {
        if (HexagonOrientation.FLAT_TOP.equals(orientation)) {
            return 1+ coordinate.getGridZ() + (coordinate.getGridX() - (coordinate.getGridX() & 1)) / 2;
        } else {
            return 1+coordinate.getGridZ();
        }
    }

    /**
     * Calculates the offset column based on a CubeCoordinate.
     *
     * @param coordinate  a cube coordinate
     * @param orientation orientation
     *
     * @return offset column or x-value
     */
    public static int convertCubeCoordinateToOffsetColumn(CubeCoordinate coordinate, HexagonOrientation orientation) {
        if (HexagonOrientation.FLAT_TOP.equals(orientation)) {
            return 1+ coordinate.getGridX();
        } else {
            return 1+ coordinate.getGridX() + (coordinate.getGridZ() - (coordinate.getGridZ() & 1)) / 2;
        }
    }
}
