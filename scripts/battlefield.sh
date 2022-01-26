#!/bin/sh
docker build -t mekstrike-battlefield -f ../battlefield/Dockerfile ../battlefield
kubectl delete --ignore-not-found=true -f ../k8s/battlefield.yaml
kubectl apply -f ../k8s/battlefield.yaml
