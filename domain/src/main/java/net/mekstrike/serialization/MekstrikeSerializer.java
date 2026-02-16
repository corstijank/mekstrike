package net.mekstrike.serialization;

import java.io.IOException;
import java.lang.reflect.InvocationTargetException;
import java.util.List;
import java.util.StringJoiner;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.protobuf.Message.Builder;
import com.fasterxml.jackson.databind.JavaType;
import com.google.protobuf.MessageOrBuilder;
import com.google.protobuf.util.JsonFormat;

import io.dapr.client.ObjectSerializer;
import io.dapr.serializer.DaprObjectSerializer;
import io.dapr.utils.TypeRef;

public class MekstrikeSerializer extends ObjectSerializer implements DaprObjectSerializer {
    private static final Logger logger = LoggerFactory.getLogger(MekstrikeSerializer.class);

    private JsonFormat.Printer printer = JsonFormat.printer();
    private JsonFormat.Parser parser = JsonFormat.parser();

    @Override
    public byte[] serialize(Object o) throws IOException {
        return switch (o) {
            case MessageOrBuilder m -> printer.print(m).getBytes();
            case List<?> list -> serializeList(list);
            default -> super.serialize(o);
        };
    }

    private byte[] serializeList(List<?> list) throws IOException {
        var joiner = new StringJoiner(",", "[", "]");
        for (Object obj : list) {
            if (obj instanceof MessageOrBuilder m) {
                joiner.add(printer.print(m));
            } else {
                joiner.add(new String(super.serialize(obj)));
            }
        }
        return joiner.toString().getBytes();
    }

    public <T> T deserialize(byte[] content, TypeRef<T> type) throws IOException {
        return deserialize(content, OBJECT_MAPPER.constructType(type.getType()));
    }

    public <T> T deserialize(byte[] content, Class<T> clazz) throws IOException {
        return deserialize(content, OBJECT_MAPPER.constructType(clazz));
    }

    @SuppressWarnings("unchecked")
    private <T> T deserialize(byte[] content, JavaType javaType) throws IOException {
        if (javaType.isPrimitive()) {
            return deserializePrimitives(content, javaType);
        }
        try {
            var newBuilderMethod = javaType.getRawClass().getMethod("newBuilder");
            var builderObj = newBuilderMethod.invoke(null);

            if (builderObj instanceof Builder builder) {
                parser.merge(new String(content), builder);
                return (T) builder.build();
            }

            throw new IllegalStateException("newBuilder() did not return a protobuf Builder");
        } catch (NoSuchMethodException _) {
            // Fallback to jackson mapper if this isn't a protobuf builder
            return OBJECT_MAPPER.readValue(content, javaType);
        } catch (SecurityException | IllegalAccessException
                | IllegalArgumentException | InvocationTargetException e) {
            logger.error("Error deserializing", e);
            throw new IllegalStateException("Unable to deserialize", e);
        }
    }

    @Override
    public String getContentType() {
        return "application/json";
    }

    @SuppressWarnings("unchecked")
    private static <T> T deserializePrimitives(byte[] content, JavaType javaType) throws IOException {
        if (content == null || content.length == 0) {
            if (javaType.hasRawClass(boolean.class)) {
                return (T) Boolean.FALSE;
            }

            if (javaType.hasRawClass(byte.class)) {
                return (T) Byte.valueOf((byte) 0);
            }

            if (javaType.hasRawClass(short.class)) {
                return (T) Short.valueOf((short) 0);
            }

            if (javaType.hasRawClass(int.class)) {
                return (T) Integer.valueOf(0);
            }

            if (javaType.hasRawClass(long.class)) {
                return (T) Long.valueOf(0L);
            }

            if (javaType.hasRawClass(float.class)) {
                return (T) Float.valueOf(0);
            }

            if (javaType.hasRawClass(double.class)) {
                return (T) Double.valueOf(0);
            }

            if (javaType.hasRawClass(char.class)) {
                return (T) Character.valueOf(Character.MIN_VALUE);
            }
            return null;
        }
        return OBJECT_MAPPER.readValue(content, javaType);
    }
}
