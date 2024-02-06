package net.mekstrike.battlefield;

import java.util.List;

public class UnitStats {
    private String name;
    private String model;
    private int pointvalue;
    private String type;
    private int size;
    private String movement;
    private String role;
    private int shortdmg;
    private int meddmg;
    private int longdmg;
    private int ovhdmg;
    private int struct;
    private List<String> specials;
    private String image;

    public String getName() {
        return this.name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getModel() {
        return this.model;
    }

    public void setModel(String model) {
        this.model = model;
    }

    public int getPointvalue() {
        return this.pointvalue;
    }

    public void setPointvalue(int pointvalue) {
        this.pointvalue = pointvalue;
    }

    public String getType() {
        return this.type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public int getSize() {
        return this.size;
    }

    public void setSize(int size) {
        this.size = size;
    }

    public String getMovement() {
        return this.movement;
    }

    public void setMovement(String movement) {
        this.movement = movement;
    }

    public String getRole() {
        return this.role;
    }

    public void setRole(String role) {
        this.role = role;
    }

    public int getShortdmg() {
        return this.shortdmg;
    }

    public void setShortdmg(int shortdmg) {
        this.shortdmg = shortdmg;
    }

    public int getMeddmg() {
        return this.meddmg;
    }

    public void setMeddmg(int meddmg) {
        this.meddmg = meddmg;
    }

    public int getLongdmg() {
        return this.longdmg;
    }

    public void setLongdmg(int longdmg) {
        this.longdmg = longdmg;
    }

    public int getOvhdmg() {
        return this.ovhdmg;
    }

    public void setOvhdmg(int ovhdmg) {
        this.ovhdmg = ovhdmg;
    }

    public int getStruct() {
        return this.struct;
    }

    public void setStruct(int struct) {
        this.struct = struct;
    }

    public List<String> getSpecials() {
        return this.specials;
    }

    public void setSpecials(List<String> specials) {
        this.specials = specials;
    }

    public String getImage() {
        return this.image;
    }

    public void setImage(String image) {
        this.image = image;
    }
}