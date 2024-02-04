#!/bin/bash
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}Adding HELM repositories${NC}" 
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

echo -e "${GREEN}Installing CertManager${NC}" 
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.12.7/cert-manager.yaml
echo -e "${GREEN}Waiting 30s for certmanager  to initialize${NC}" 
sleep 30

echo -e "${GREEN}Installing Opentelemetry Operator${NC}" 
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml

echo -e "${GREEN}Installing Jaeger Operator${NC}" 
kubectl create namespace observability
kubectl create namespace monitoring

kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.52.0/jaeger-operator.yaml -n observability # <2>
echo -e "${GREEN}Waiting 30s for otel and jaeger operator to initialize${NC}" 
sleep 30

echo -e "${GREEN}Deplaying Otel collector${NC}" 
kubectl apply -f open-telemetry-collector.yaml

echo -e "${GREEN}Installing Jaeger${NC}" 
kubectl apply -f jaeger.yaml

echo -e "${GREEN}Installing Helm Apps and operators, Dapr and Redis${NC}" 
helm upgrade --install dapr dapr/dapr --version=1.12.0 --namespace dapr-system --create-namespace
helm upgrade --install --values=./redis-values.yaml redis bitnami/redis --namespace redis --create-namespace 

echo -e "${GREEN}Waiting 30s for everything else to initialize${NC}"
sleep 30
