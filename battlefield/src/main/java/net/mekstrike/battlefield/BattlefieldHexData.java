package net.mekstrike.battlefield;

import org.hexworks.mixite.core.api.defaults.DefaultSatelliteData;

public class BattlefieldHexData extends DefaultSatelliteData {
    private boolean occupied;

    public boolean isOccupied() {
        return occupied;
    }

    public void setOccupied(boolean occupied) {
        this.occupied = occupied;
    }
}
