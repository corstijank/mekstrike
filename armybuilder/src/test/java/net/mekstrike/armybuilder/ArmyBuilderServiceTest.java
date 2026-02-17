package net.mekstrike.armybuilder;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import com.google.protobuf.util.JsonFormat;

import io.dapr.client.DaprClient;
import io.dapr.client.domain.HttpExtension;
import net.mekstrike.armybuilder.ArmyBuilderProtos.ArmyRequest;
import net.mekstrike.domain.unit.Unit;
import reactor.core.publisher.Mono;

class ArmyBuilderServiceTest {

    private DaprClient mockClient;
    private ArmyBuilderService service;
    private final JsonFormat.Printer printer = JsonFormat.printer();

    @BeforeEach
    void setUp() {
        mockClient = mock(DaprClient.class);
        service = new ArmyBuilderService(mockClient);
    }

    @Test
    void createArmyReturnsEmptyResponseWhenAllCountsZero() throws Exception {
        var request = ArmyRequest.newBuilder().build();

        var response = service.createArmy(request);

        assertEquals(0, response.getUnitsCount());
    }

    @Test
    void createArmyFetchesCorrectNumberOfLightUnits() throws Exception {
        stubLibraryResponse("light");

        var request = ArmyRequest.newBuilder().setLights(2).build();
        var response = service.createArmy(request);

        assertEquals(2, response.getUnitsCount());
        verify(mockClient, times(2)).invokeMethod(
                eq("library"), eq("units/by/BM/light/random"),
                eq(""), eq(HttpExtension.GET), isNull(), eq(byte[].class));
    }

    @Test
    void createArmyFetchesAllWeightClasses() throws Exception {
        stubLibraryResponse("light");
        stubLibraryResponse("medium");
        stubLibraryResponse("heavy");
        stubLibraryResponse("assault");

        var request = ArmyRequest.newBuilder()
                .setLights(1)
                .setMediums(1)
                .setHeavies(1)
                .setAssaults(1)
                .build();

        var response = service.createArmy(request);

        assertEquals(4, response.getUnitsCount());
        verify(mockClient).invokeMethod(eq("library"), eq("units/by/BM/light/random"),
                any(), any(), any(), eq(byte[].class));
        verify(mockClient).invokeMethod(eq("library"), eq("units/by/BM/medium/random"),
                any(), any(), any(), eq(byte[].class));
        verify(mockClient).invokeMethod(eq("library"), eq("units/by/BM/heavy/random"),
                any(), any(), any(), eq(byte[].class));
        verify(mockClient).invokeMethod(eq("library"), eq("units/by/BM/assault/random"),
                any(), any(), any(), eq(byte[].class));
    }

    @Test
    void createArmyParsesUnitStatsFromLibraryResponse() throws Exception {
        var stats = Unit.Stats.newBuilder()
                .setName("Atlas")
                .setModel("AS7-D")
                .setPointvalue(52)
                .setType("BM")
                .build();

        when(mockClient.invokeMethod(eq("library"), contains("heavy"),
                any(), any(), any(), eq(byte[].class)))
                .thenReturn(Mono.just(printer.print(stats).getBytes()));

        var request = ArmyRequest.newBuilder().setHeavies(1).build();
        var response = service.createArmy(request);

        assertEquals(1, response.getUnitsCount());
        var unit = response.getUnits(0);
        assertEquals("Atlas", unit.getName());
        assertEquals("AS7-D", unit.getModel());
        assertEquals(52, unit.getPointvalue());
        assertEquals("BM", unit.getType());
    }

    @Test
    void createArmyReturnsMixedArmyWithCorrectTotalCount() throws Exception {
        stubLibraryResponse("light");
        stubLibraryResponse("medium");
        stubLibraryResponse("heavy");
        stubLibraryResponse("assault");

        var request = ArmyRequest.newBuilder()
                .setLights(3)
                .setMediums(2)
                .setHeavies(2)
                .setAssaults(1)
                .build();

        var response = service.createArmy(request);

        assertEquals(8, response.getUnitsCount());
    }

    private void stubLibraryResponse(String weightClass) throws Exception {
        var stats = Unit.Stats.newBuilder()
                .setName("Test " + weightClass)
                .setType("BM")
                .build();

        when(mockClient.invokeMethod(eq("library"),
                eq("units/by/BM/%s/random".formatted(weightClass)),
                eq(""), eq(HttpExtension.GET), isNull(), eq(byte[].class)))
                .thenReturn(Mono.just(printer.print(stats).getBytes()));
    }
}
