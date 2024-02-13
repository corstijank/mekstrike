# Build JAR
FROM maven:3-eclipse-temurin-21 AS BUILDER

WORKDIR /app

COPY unit /app/unit
COPY domain /app/domain

WORKDIR /app/domain
RUN mvn dependency:go-offline
RUN mvn install

WORKDIR /app/unit
RUN mvn dependency:go-offline
RUN mvn package

## Build actual image
FROM eclipse-temurin:21

WORKDIR /app

# COPY the fat jar
COPY --from=builder /app/unit/target/mekstrike-unit-exec.jar /app/unit.jar

CMD ["java","-jar","/app/unit.jar"]