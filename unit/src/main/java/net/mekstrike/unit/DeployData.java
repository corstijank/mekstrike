package net.mekstrike.unit;

import com.fasterxml.jackson.annotation.JsonProperty;

// TODO: That fancy record stuff. 
public class DeployData {
    @JsonProperty("BattlefieldID")
    private String battlefieldID;

    @JsonProperty("Owner")
    private String owner;

    @JsonProperty("Stats")
    private UnitStats stats;

    @JsonProperty("DeployLocation")
    private String deployLocation;

    public String getBattlefieldID() {
        return battlefieldID;
    }

    public void setBattlefieldID(String battlefieldID) {
        this.battlefieldID = battlefieldID;
    }

    public String getOwner() {
        return owner;
    }

    public void setOwner(String owner) {
        this.owner = owner;
    }

    public UnitStats getStats() {
        return stats;
    }

    public void setStats(UnitStats stats) {
        this.stats = stats;
    }

    public String getDeployLocation() {
        return deployLocation;
    }

    public void setDeployLocation(String deployLocation) {
        this.deployLocation = deployLocation;
    }
}
