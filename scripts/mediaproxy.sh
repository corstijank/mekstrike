#!/bin/sh
docker build -t mekstrike-mediaproxy -f ../mediaproxy/Dockerfile .. --no-cache
kubectl delete --ignore-not-found=true -f ../k8s/mediaproxy.yaml
kubectl apply -f ../k8s/mediaproxy.yaml
