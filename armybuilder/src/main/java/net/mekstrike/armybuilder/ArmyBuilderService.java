package net.mekstrike.armybuilder;

import java.io.IOException;


import com.google.protobuf.Any;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.util.JsonFormat;
import io.dapr.client.DaprClient;
import io.dapr.client.DaprClientBuilder;
import io.dapr.client.domain.HttpExtension;
import io.dapr.v1.AppCallbackGrpc;
import io.dapr.v1.CommonProtos;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import net.mekstrike.armybuilder.ArmyBuilderProtos.ArmyRequest;
import net.mekstrike.armybuilder.ArmyBuilderProtos.ArmyResponse;
import net.mekstrike.common.unit.CommonUnitProtos;

public class ArmyBuilderService extends AppCallbackGrpc.AppCallbackImplBase {

  private Server server;
  private DaprClient client;

  /**
   * Server mode: starts listening on given port.
   *
   * @param port Port to listen on.
   * @throws IOException Errors while trying to start service.
   */
  public void start(int port) throws IOException {
    this.server = ServerBuilder.forPort(port).addService(this).build().start();
    System.out.printf("Server: started listening on port %d\n", port);
    client = new DaprClientBuilder().build();
    // Now we handle ctrl+c (or any other JVM shutdown)
    Runtime.getRuntime().addShutdownHook(new Thread() {

      @Override
      public void run() {
        try {
          System.out.println("Server: shutting down gracefully ...");
          ArmyBuilderService.this.server.shutdown();
          client.close();
        } catch (Exception e) {
          e.printStackTrace();
        }
        System.out.println("Server: Bye.");
      }
    });
  }

  /**
   * Server mode: waits for shutdown trigger.
   *
   * @throws InterruptedException Propagated interrupted exception.
   */
  public void awaitTermination() throws InterruptedException {
    if (this.server != null) {
      this.server.awaitTermination();
    }
  }

  /**
   * Server mode: this is the Dapr method to receive Invoke operations via Grpc.
   *
   * @param request          Dapr envelope request,
   * @param responseObserver Dapr envelope response.
   */
  @Override
  public void onInvoke(CommonProtos.InvokeRequest request,
      StreamObserver<CommonProtos.InvokeResponse> responseObserver) {
    try {
      if ("createArmy".equals(request.getMethod())) {
        ArmyRequest armyRequest = ArmyRequest.newBuilder().mergeFrom(request.getData().getValue()).build();

        var army = createArmy(armyRequest);

        // Create an the dapr envelope response
        CommonProtos.InvokeResponse.Builder responseBuilder = CommonProtos.InvokeResponse.newBuilder();
        responseBuilder.setData(Any.pack(army));
        responseObserver.onNext(responseBuilder.build());
      }
    } catch (InvalidProtocolBufferException e) {
      System.out.println("Error merging messages!");
      e.printStackTrace();
    } finally {
      responseObserver.onCompleted();
    }
  }

  public ArmyResponse createArmy(ArmyRequest req) throws InvalidProtocolBufferException {
    var parser = JsonFormat.parser();
    var result = ArmyResponse.newBuilder();
    for (int i = 0; i < req.getLights(); i++) {
      var unit = CommonUnitProtos.UnitStats.newBuilder();
      byte[] response = client
          .invokeMethod("library", "units/by/BM/light/random", "", HttpExtension.GET, null, byte[].class).block();

      parser.merge(new String(response), unit);
      result.addUnits(unit);

    }
    for (int i = 0; i < req.getMediums(); i++) {
      var unit = CommonUnitProtos.UnitStats.newBuilder();
      byte[] response = client
          .invokeMethod("library", "units/by/BM/medium/random", "", HttpExtension.GET, null, byte[].class).block();
      parser.merge(new String(response), unit);
      result.addUnits(unit);
    }
    for (int i = 0; i < req.getHeavies(); i++) {
      var unit = CommonUnitProtos.UnitStats.newBuilder();
      byte[] response = client
          .invokeMethod("library", "units/by/BM/heavy/random", "", HttpExtension.GET, null, byte[].class).block();
      parser.merge(new String(response), unit);
      result.addUnits(unit);
    }
    for (int i = 0; i < req.getAssaults(); i++) {
      var unit = CommonUnitProtos.UnitStats.newBuilder();
      byte[] response = client
          .invokeMethod("library", "units/by/BM/assault/random", "", HttpExtension.GET, null, byte[].class).block();
      parser.merge(new String(response), unit);
      result.addUnits(unit);
    }

    return result.build();
  }
}
