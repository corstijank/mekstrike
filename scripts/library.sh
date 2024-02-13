#!/bin/sh
nerdctl --namespace k8s.io build -t mekstrike-library -f ../library/Dockerfile ..
kubectl delete --ignore-not-found=true -f ../k8s/library.yaml
kubectl apply -f ../k8s/library.yaml
