# Build JAR
FROM maven:3-eclipse-temurin-21 AS BUILDER

WORKDIR /app

COPY battlefield/pom.xml .

RUN mvn dependency:go-offline

COPY battlefield/src/ ./src/

RUN mvn package

## Build actual image
FROM eclipse-temurin:21

WORKDIR /app

# COPY the fat jar
COPY --from=builder /app/target/mekstrike-battlefield-exec.jar /app/battlefield.jar

CMD ["java","-jar","/app/battlefield.jar"]