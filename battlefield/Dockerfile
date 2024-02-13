# Build JAR
FROM maven:3-eclipse-temurin-21 AS BUILDER

WORKDIR /app

COPY battlefield /app/battlefield
COPY domain /app/domain

WORKDIR /app/domain
RUN mvn dependency:go-offline
RUN mvn install

WORKDIR /app/battlefield
RUN mvn dependency:go-offline
RUN mvn package

## Build actual image
FROM eclipse-temurin:21

WORKDIR /app

# COPY the fat jar
COPY --from=builder /app/battlefield/target/mekstrike-battlefield-exec.jar /app/battlefield.jar

CMD ["java","-jar","/app/battlefield.jar"]