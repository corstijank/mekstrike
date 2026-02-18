# Build JAR
FROM maven:3-eclipse-temurin-25 AS BUILDER

COPY pom.xml /app/pom.xml
COPY domain /app/domain
COPY armybuilder /app/armybuilder

WORKDIR /app/armybuilder
RUN mvn dependency:go-offline
RUN mvn package

## Build actual image
FROM eclipse-temurin:25

WORKDIR /app

# COPY jar
COPY --from=builder /app/armybuilder/target/armybuilder-1.0.0-SNAPSHOT-jar-with-dependencies.jar /app/armybuilder.jar

CMD ["java","-jar","/app/armybuilder.jar"]
