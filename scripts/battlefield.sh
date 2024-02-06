#!/bin/sh
nerdctl --namespace k8s.io build -t mekstrike-battlefield -f ../battlefield/Dockerfile ..
kubectl delete --ignore-not-found=true -f ../k8s/battlefield.yaml
kubectl apply -f ../k8s/battlefield.yaml
