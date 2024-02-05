package net.mekstrike.unit;

public class UnitData {
    private UnitLocation location;

    private UnitStats stats;

    private String owner;

    private boolean active;

    public UnitData() {

    }

    public UnitData(UnitLocation location, UnitStats stats, String owner, boolean active) {
        this.location = location;
        this.stats = stats;
        this.owner = owner;
        this.active = active;
    }

    public UnitLocation getLocation() {
        return location;
    }

    public void setLocation(UnitLocation location) {
        this.location = location;
    }

    public UnitStats getStats() {
        return stats;
    }

    public void setStats(UnitStats stats) {
        this.stats = stats;
    }

    public String getOwner() {
        return owner;
    }

    public void setOwner(String owner) {
        this.owner = owner;
    }

    public boolean isActive() {
        return active;
    }

    public void setActive(boolean active) {
        this.active = active;
    }
}
