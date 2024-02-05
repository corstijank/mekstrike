#!/bin/sh
nerdctl --namespace k8s.io build -t mekstrike-armybuilder -f ../armybuilder/Dockerfile ../armybuilder
kubectl delete --ignore-not-found=true -f ../k8s/armybuilder.yaml
kubectl apply -f ../k8s/armybuilder.yaml
