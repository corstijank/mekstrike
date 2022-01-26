# Mekstrike

## Running / trying

### Requirements

- Kubernetes with an ability to build images on the nodes. For simplicity sake; docker-desktop or minikube (untested)

### Quickstart

```sh
# Set up a kubernetes platorm
./k8s/platform/setup.sh
# This will be downloading bits in the background, meanwhile we'll go ahead and build some images
# Build images
./bldImgs.sh
# Deploy application, make sure dapr is up and running before running this
./k8s/deploy.sh
# Start a new game
curl -d '{"PlayerName":"HelloThisIsDog"}' -H "Content-Type: application/json" -X POST http://localhost/mekstrike/api/gamemaster/games
# Or use the excellent HTTPie:
# http post http://localhost/mekstrike/api/gamemaster/games PlayerName=ThisIsDog
```

### Useful CLI snippets

```sh
# Open Traefik Dashboard (Open http://localhost:9000/dashboard)
kubectl port-forward $(kubectl get pods --namespace traefik --selector "app.kubernetes.io/name=traefik" --output=name) 9000:9000 --namespace traefik

# Get Elastic password:
kubectl get secret elasticsearch-es-elastic-user -o=jsonpath='{.data.elastic}' --namespace monitoring| base64 --decode; echo
```

## Developing: requirements

- GoLang 1.16 (or later)
- DotNet 6.0 (or later)
- Java 16 (or later)
- Maven (something recent)

## Components

desired:
![overview](overview.svg)

### library

- golang
- rest endpoints for querying unit data

### armybuilder

- java
- grpc service for creating an army
  - Takes input parameters for generating a force
  - Generate force - query library for units
  - creates and returns unit URI's
  
### gamemaster

- Go
- REST
  - players
  - the battlefield
  - involved units
  - rounds / phases

### battlefield

- .NET
- Actor
  - Given board x, can a unit move from x to y?
  - etc

### unit

- Java
- Actor
  - Move here
  - Fire here
  - handle IncomingFire
  - etc
  - knows which map its on, and its location

### fireControlSystem

- Python for shits and giggles?
  - will translate *shots* to *incoming fire*, given RNG and all kinds of considerations;
  - Client -> Unit(fire) -> FireControlSystem(verifyShot) -> target(processIncoming)
  
### API Gateway

- A simple API gateway so client does not need knowledge of backend services?

### UI

- Svelte webapp
- CSS: https://terminalcss.xyz/ 
