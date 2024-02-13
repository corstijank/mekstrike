package net.mekstrike.armybuilder;

import java.io.IOException;

public class Program {
    public static void main(String[] args) throws IOException, InterruptedException {
        new Program().start();
    }

    public void start() throws IOException, InterruptedException {
        final ArmyBuilderService service = new ArmyBuilderService();
        service.start(9000);
        service.awaitTermination();
    }
}
