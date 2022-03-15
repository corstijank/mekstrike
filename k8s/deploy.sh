#!/bin/bash
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}Deploying Namespace${NC}"
kubectl apply -f namespace.yaml

echo -e "${GREEN}Deploying Instrumentation${NC}"
kubectl apply -f mekstrike-instrumentation.yaml

echo -e "${GREEN}Deploying Dapr configs and components${NC}"
kubectl apply -f mekstrike-config.yaml
kubectl apply -f library-store.yaml
kubectl apply -f battlefield-store.yaml

echo -e "${GREEN}Deploying Backend services${NC}"
kubectl apply -f library.yaml
kubectl apply -f armybuilder.yaml
kubectl apply -f gamemaster.yaml
kubectl apply -f battlefield.yaml
kubectl apply -f unit.yaml

echo -e "${GREEN}Deploying Frontend${NC}"
kubectl apply -f mediaproxy.yaml
kubectl apply -f ui.yaml

echo -e "${GREEN}Waiting for pods to initialize before applying importer${NC}"
sleep 30
kubectl apply -f library-importer.yaml

echo -e "${GREEN}Waiting for importer to complete${NC}"
kubectl wait --for=condition=complete --timeout=240s job/library-importer  --namespace mekstrike