#!/bin/sh
docker build -t mekstrike-ui -f ../ui/Dockerfile ../ui/ --no-cache
kubectl delete --ignore-not-found=true -f ../k8s/ui.yaml
kubectl apply -f ../k8s/ui.yaml
