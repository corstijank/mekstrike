package net.mekstrike.armybuilder;

import io.dapr.client.DaprClientBuilder;
import io.grpc.ServerBuilder;
import java.io.IOException;

public class Program {
  public static void main(String[] args) throws IOException, InterruptedException {
    new Program().start();
  }

  public void start() throws IOException, InterruptedException {
    var client = new DaprClientBuilder().build();
    var service = new ArmyBuilderService(client);
    var server = ServerBuilder.forPort(9000).addService(service).build().start();
    System.out.printf("Server: started listening on port %d%n", 9000);

    Runtime.getRuntime()
        .addShutdownHook(
            new Thread(
                () -> {
                  try {
                    System.out.println("Server: shutting down gracefully ...");
                    server.shutdown();
                    client.close();
                  } catch (Exception e) {
                    e.printStackTrace();
                  }
                  System.out.println("Server: Bye.");
                }));

    server.awaitTermination();
  }
}
