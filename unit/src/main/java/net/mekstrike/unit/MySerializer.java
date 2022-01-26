package net.mekstrike.unit;

import java.io.IOException;

import io.dapr.serializer.DaprObjectSerializer;
import io.dapr.utils.TypeRef;

public class MySerializer implements DaprObjectSerializer {

    @Override
    public byte[] serialize(Object o) throws IOException {
        return null;
    }

    @Override
    public <T> T deserialize(byte[] data, TypeRef<T> type) throws IOException {
        System.out.println(new String(data));
        return null;
    }

    @Override
    public String getContentType() {
        return "application/json";
    }

}
