# Build JAR
FROM openjdk:17-jdk AS BUILDER

WORKDIR /app

COPY . .
WORKDIR /app/armybuilder

RUN ./mvnw package

## Build actual image
FROM openjdk:17-jdk

WORKDIR /app

# COPY jar
COPY --from=builder /app/armybuilder/target/armybuilder-1.0.0-SNAPSHOT-jar-with-dependencies.jar /app/armybuilder.jar

CMD ["java","-jar","/app/armybuilder.jar"]