# Build JAR
FROM openjdk:11-jdk AS BUILDER

WORKDIR /app

COPY . .
WORKDIR /app/unit

# Build the jar
RUN ./mvnw package

## Build actual image
FROM openjdk:11-jdk

WORKDIR /app

# COPY the fat jar
COPY --from=builder /app/unit/target/mekstrike-unit-exec.jar /app/unit.jar

CMD ["java","-jar","/app/unit.jar"]