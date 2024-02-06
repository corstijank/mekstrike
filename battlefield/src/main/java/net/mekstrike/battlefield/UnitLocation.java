package net.mekstrike.battlefield;

public class UnitLocation {
    private String battlefieldID;
    private Coordinates position;
    private int heading;

    public UnitLocation() {
        // Default constructor to facilitate deserializing by Dapr 
    }

    public UnitLocation(String battlefieldID, Coordinates position, int heading) {
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

    public Coordinates getPosition() {
        return this.position;
    }

    public void setPosition(Coordinates position) {
        this.position = position;
    }

    public int getHeading() {
        return this.heading;
    }

    public void setHeading(int heading) {
        this.heading = heading;
    }

}
