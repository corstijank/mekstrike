package net.mekstrike.unit;

import com.fasterxml.jackson.annotation.JsonProperty;

public class CellRef {
    @JsonProperty("Col")
    private int col;

    @JsonProperty("Row")
    private int row;

    public CellRef(){
        // Default construction to allow deserialization by Dapr
    }

    public CellRef( int col, int row) {
        this.col = col;
        this.row = row;
    }

    public int getCol() {
        return this.col;
    }

    public void setCol(int col) {
        this.col = col;
    }

    public int getRow() {
        return this.row;
    }

    public void setRow(int row) {
        this.row = row;
    }

}
