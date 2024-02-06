package net.mekstrike.battlefield;

import io.dapr.actors.ActorId;
import io.dapr.actors.runtime.AbstractActor;
import io.dapr.actors.runtime.ActorRuntimeContext;
import org.hexworks.mixite.core.api.*;
import reactor.core.publisher.Mono;

import java.util.ArrayList;
import java.util.List;

public class BattlefieldImpl extends AbstractActor implements Battlefield{
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
    }

    @Override
    public List<Cell> getBoardCells() {
        var result = new ArrayList<Cell>();
        for (Hexagon<BattlefieldHexData> hexagon : grid.getHexagons()) {
            var c = hexagon.getCubeCoordinate();
            var x = convertCubeCoordinateToOffsetColumn(c, HexagonOrientation.FLAT_TOP);
            var y = convertCubeCoordinateToOffsetRow(c, HexagonOrientation.FLAT_TOP);
            result.add(new Cell(new Coordinates(x, y), ""));
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
    public boolean isCellBlocked(Coordinates cellRef) {
        var result=false;
        var x = CoordinateConverter.convertOffsetCoordinatesToCubeX(cellRef.getCol(),cellRef.getRow(),HexagonOrientation.FLAT_TOP);
        var z = CoordinateConverter.convertOffsetCoordinatesToCubeZ(cellRef.getCol(),cellRef.getRow(),HexagonOrientation.FLAT_TOP);

        CubeCoordinate c = CubeCoordinate.fromCoordinates(x,z);
        var maybeHex = grid.getByCubeCoordinate(c);
        if(maybeHex.isPresent()){
            var h = maybeHex.get();
            if(!h.getSatelliteData().isPresent()){
                h.setSatelliteData(new BattlefieldHexData());
            }
            result= h.getSatelliteData().get().isOccupied();
        }
        return result;
    }

    @Override
    public void blockCell(Coordinates cellRef) {
        var x = CoordinateConverter.convertOffsetCoordinatesToCubeX(cellRef.getCol(),cellRef.getRow(),HexagonOrientation.FLAT_TOP);
        var z = CoordinateConverter.convertOffsetCoordinatesToCubeZ(cellRef.getCol(),cellRef.getRow(),HexagonOrientation.FLAT_TOP);

        CubeCoordinate c = CubeCoordinate.fromCoordinates(x,z);
        var maybeHex = grid.getByCubeCoordinate(c);
        if(maybeHex.isPresent()){
            var h = maybeHex.get();
            if(!h.getSatelliteData().isPresent()){
                h.setSatelliteData(new BattlefieldHexData());
            }
            h.getSatelliteData().get().setOccupied(true);
        }
    }

    @Override
    public List<Coordinates> getMovementOptions(UnitData unit) {
        var result = new ArrayList<Coordinates>();

        var x = CoordinateConverter.convertOffsetCoordinatesToCubeX(unit.getLocation().getPosition().getCol(),unit.getLocation().getPosition().getRow(),HexagonOrientation.FLAT_TOP);
        var z = CoordinateConverter.convertOffsetCoordinatesToCubeZ(unit.getLocation().getPosition().getCol(),unit.getLocation().getPosition().getRow(),HexagonOrientation.FLAT_TOP);

        CubeCoordinate c = CubeCoordinate.fromCoordinates(x,z);
        var maybeHex = grid.getByCubeCoordinate(c);
        if(maybeHex.isPresent()){
            var h = maybeHex.get();
            if(!h.getSatelliteData().isPresent()){
                var calculator = new HexagonalGridBuilder<BattlefieldHexData>().buildCalculatorFor(grid);
                var movement = Integer.parseInt(unit.getStats().getMovement().split("\"")[0])/2;
                var options = calculator.calculateMovementRangeFrom(h,movement);
                for (Hexagon<BattlefieldHexData> hexagon : options) {
                    var oc = hexagon.getCubeCoordinate();
                    var ox = convertCubeCoordinateToOffsetColumn(oc, HexagonOrientation.FLAT_TOP);
                    var oy = convertCubeCoordinateToOffsetRow(oc, HexagonOrientation.FLAT_TOP);
                    result.add(new Coordinates(ox, oy));
                }
            }
        }
        return result;
    }


    public Mono<Void>saveBattlefield(){
        return getActorStateManager().set("grid", grid);
    }

    @Override
    protected Mono<Void> onActivate() {
        // grid = getActorStateManager().get("grid",HexagonalGrid.class).block();
        return Mono.empty();
    }

    @Override
    protected Mono<Void> onDeactivate() {
        return saveBattlefield();
    }


    /**
     * Calculates the offset row based on a CubeCoordinate.
     *
     * @param coordinate a cube coordinate
     * @param orientation orientation
     *
     * @return offset row or y-value
     */
    public static int convertCubeCoordinateToOffsetRow(CubeCoordinate coordinate, HexagonOrientation orientation) {
        if(HexagonOrientation.FLAT_TOP.equals(orientation)) {
            return coordinate.getGridZ() + (coordinate.getGridX() - (coordinate.getGridX() & 1)) / 2;
        } else {
            return coordinate.getGridZ();
        }
    }

    /**
     * Calculates the offset column based on a CubeCoordinate.
     *
     * @param coordinate a cube coordinate
     * @param orientation orientation
     *
     * @return offset column or x-value
     */
    public static int convertCubeCoordinateToOffsetColumn(CubeCoordinate coordinate, HexagonOrientation orientation) {
        if(HexagonOrientation.FLAT_TOP.equals(orientation)) {
            return coordinate.getGridX();
        } else {
            return coordinate.getGridX() + (coordinate.getGridZ() - (coordinate.getGridZ() & 1)) / 2;
        }
    }
}
