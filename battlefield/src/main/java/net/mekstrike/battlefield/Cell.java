package net.mekstrike.battlefield;

public class Cell {
    private Coordinates coordinates;
    private String type;

    public Cell(Coordinates c, String t){
        this.coordinates=c;
        this.type=t;
    }

    public Coordinates getCoordinates() {
        return coordinates;
    }

    public String getType() {
        return type;
    }
}
