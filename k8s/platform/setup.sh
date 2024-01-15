#!/bin/bash
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}Adding HELM repositories${NC}" 
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

echo -e "${GREEN}Installing CertManager and ECK Operator${NC}" 
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml
kubectl apply -f crds.yaml
kubectl apply -f operator.yaml
echo -e "${GREEN}Waiting 30s for eck operator to initialize${NC}" 
sleep 30

echo -e "${GREEN}Installing Opentelemetry Operator${NC}" 
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml

echo -e "${GREEN}Waiting 30s for otel to initialize${NC}" 
sleep 30

echo -e "${GREEN}Deplaying elastic components${NC}" 
kubectl apply -f monitoring-namespace.yaml
kubectl apply -f elastic.yaml
kubectl apply -f kibana.yaml
kubectl apply -f apm.yaml
kubectl apply -f beats.yaml

echo -e "${GREEN}Deplaying Otel collector${NC}" 
kubectl apply -f open-telemetry-collector.yaml

echo -e "${GREEN}Installing Helm Apps and operators, Dapr, Traefik and Redis${NC}" 
helm upgrade --install dapr dapr/dapr --version=1.12.0 --namespace dapr-system --create-namespace
helm upgrade --install --values=./redis-values.yaml redis bitnami/redis --namespace redis --create-namespace 


echo -e "${GREEN}Waiting 30s for everything else to initialize${NC}"
sleep 30
