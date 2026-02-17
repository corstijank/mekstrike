package net.mekstrike.armybuilder;

import com.google.protobuf.Any;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.util.JsonFormat;
import io.dapr.client.DaprClient;
import io.dapr.client.domain.HttpExtension;
import io.dapr.v1.AppCallbackGrpc;
import io.dapr.v1.CommonProtos;
import io.grpc.stub.StreamObserver;
import net.mekstrike.armybuilder.ArmyBuilderProtos.ArmyRequest;
import net.mekstrike.armybuilder.ArmyBuilderProtos.ArmyResponse;
import net.mekstrike.domain.unit.Unit;

public class ArmyBuilderService extends AppCallbackGrpc.AppCallbackImplBase {

  private final DaprClient client;

  public ArmyBuilderService(DaprClient client) {
    this.client = client;
  }

  @Override
  public void onInvoke(CommonProtos.InvokeRequest request,
      StreamObserver<CommonProtos.InvokeResponse> responseObserver) {
    try {
      if ("createArmy".equals(request.getMethod())) {
        ArmyRequest armyRequest = ArmyRequest.newBuilder().mergeFrom(request.getData().getValue()).build();

        var army = createArmy(armyRequest);

        var responseBuilder = CommonProtos.InvokeResponse.newBuilder();
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

    fetchUnits(parser, result, req.getLights(), "light");
    fetchUnits(parser, result, req.getMediums(), "medium");
    fetchUnits(parser, result, req.getHeavies(), "heavy");
    fetchUnits(parser, result, req.getAssaults(), "assault");

    return result.build();
  }

  private void fetchUnits(JsonFormat.Parser parser, ArmyResponse.Builder result,
      int count, String weightClass) throws InvalidProtocolBufferException {
    for (int i = 0; i < count; i++) {
      var unit = Unit.Stats.newBuilder();
      byte[] response = client
          .invokeMethod("library", "units/by/BM/%s/random".formatted(weightClass),
              "", HttpExtension.GET, null, byte[].class)
          .block();
      parser.merge(new String(response), unit);
      result.addUnits(unit);
    }
  }
}
