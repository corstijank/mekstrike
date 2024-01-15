#!/bin/sh
nerdctl --namespace k8s.io build -t mekstrike-gamemaster -f ../gamemaster/Dockerfile ..
kubectl delete --ignore-not-found=true -f ../k8s/gamemaster.yaml
kubectl apply -f ../k8s/gamemaster.yaml
