package net.mekstrike.unit;

public class UnitLocation {
    private String battlefieldID;
    private CellRef position;
    private int heading;

    public UnitLocation() {
        // Default constructor to facilitate deserializing by Dapr 
    }

    public UnitLocation(String battlefieldID, CellRef position, int heading) {
        this.battlefieldID = battlefieldID;
        this.position = position;
        this.heading = heading;
    }

    public String getBattlefieldID() {
        return this.battlefieldID;
    }

    public void setBattlefieldID(String battlefieldID) {
        this.battlefieldID = battlefieldID;
    }

    public CellRef getPosition() {
        return this.position;
    }

    public void setPosition(CellRef position) {
        this.position = position;
    }

    public int getHeading() {
        return this.heading;
    }

    public void setHeading(int heading) {
        this.heading = heading;
    }

}
