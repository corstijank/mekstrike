#!/bin/sh
nerdctl --namespace k8s.io build -t mekstrike-unit -f ../unit/Dockerfile ..
kubectl delete --ignore-not-found=true -f ../k8s/unit.yaml
kubectl apply -f ../k8s/unit.yaml