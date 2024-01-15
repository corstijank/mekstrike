#!/bin/sh
nerdctl --namespace k8s.io build -t mekstrike-ui -f ../ui/Dockerfile ../ui/ --no-cache
kubectl delete --ignore-not-found=true -f ../k8s/ui.yaml
kubectl apply -f ../k8s/ui.yaml
