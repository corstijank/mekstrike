package net.mekstrike.serialization;

import static org.junit.jupiter.api.Assertions.*;

import java.io.IOException;
import java.util.List;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import io.dapr.utils.TypeRef;
import net.mekstrike.domain.battlefield.Battlefield;

class MekstrikeSerializerTest {

    private MekstrikeSerializer serializer;

    @BeforeEach
    void setUp() {
        serializer = new MekstrikeSerializer();
    }

    @Test
    void serializesProtobufMessageToJson() throws IOException {
        var coords = Battlefield.Coordinates.newBuilder()
                .setX(3)
                .setY(7)
                .build();

        var result = new String(serializer.serialize(coords));

        assertTrue(result.contains("\"x\": 3"));
        assertTrue(result.contains("\"y\": 7"));
    }

    @Test
    void serializesListOfProtobufMessagesToJsonArray() throws IOException {
        var c1 = Battlefield.Coordinates.newBuilder().setX(1).setY(2).build();
        var c2 = Battlefield.Coordinates.newBuilder().setX(3).setY(4).build();

        var result = new String(serializer.serialize(List.of(c1, c2)));

        assertTrue(result.startsWith("["));
        assertTrue(result.endsWith("]"));
        assertTrue(result.contains("\"x\": 1"));
        assertTrue(result.contains("\"x\": 3"));
    }

    @Test
    void serializesPlainObjectViaJacksonFallback() throws IOException {
        var result = new String(serializer.serialize("hello"));

        assertEquals("\"hello\"", result);
    }

    @Test
    void deserializesProtobufMessageFromJson() throws IOException {
        var json = "{\"x\": 5, \"y\": 10}";

        var coords = serializer.deserialize(json.getBytes(), Battlefield.Coordinates.class);

        assertEquals(5, coords.getX());
        assertEquals(10, coords.getY());
    }

    @Test
    void deserializesProtobufMessageViaTypeRef() throws IOException {
        var json = "{\"x\": 5, \"y\": 10}";

        var coords = serializer.deserialize(json.getBytes(), TypeRef.get(Battlefield.Coordinates.class));

        assertEquals(5, coords.getX());
        assertEquals(10, coords.getY());
    }

    @Test
    void deserializesNonProtobufClassViaJacksonFallback() throws IOException {
        var json = "\"hello\"";

        var result = serializer.deserialize(json.getBytes(), String.class);

        assertEquals("hello", result);
    }

    @Test
    void deserializesEmptyContentToDefaultPrimitives() throws IOException {
        assertEquals(false, serializer.deserialize(new byte[0], boolean.class));
        assertEquals((byte) 0, serializer.deserialize(new byte[0], byte.class));
        assertEquals((short) 0, serializer.deserialize(new byte[0], short.class));
        assertEquals(0, serializer.deserialize(new byte[0], int.class));
        assertEquals(0L, serializer.deserialize(new byte[0], long.class));
        assertEquals(0f, serializer.deserialize(new byte[0], float.class));
        assertEquals(0d, serializer.deserialize(new byte[0], double.class));
        assertEquals(Character.MIN_VALUE, serializer.deserialize(new byte[0], char.class));
    }

    @Test
    void deserializesNullContentToDefaultPrimitives() throws IOException {
        assertEquals(false, serializer.deserialize(null, boolean.class));
        assertEquals(0, serializer.deserialize(null, int.class));
    }

    @Test
    void roundTripsProtobufMessage() throws IOException {
        var original = Battlefield.Cell.newBuilder()
                .setCoordinates(Battlefield.Coordinates.newBuilder().setX(4).setY(8))
                .setType("forest")
                .build();

        var bytes = serializer.serialize(original);
        var restored = serializer.deserialize(bytes, Battlefield.Cell.class);

        assertEquals(original, restored);
    }

    @Test
    void contentTypeIsJson() {
        assertEquals("application/json", serializer.getContentType());
    }
}
