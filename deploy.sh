#!/bin/bash
GREEN='\033[0;32m'
NC='\033[0m' # No Color

build_image() {
    local namespace="k8s.io"
    local image_name="$1"
    local dockerfile_path="$image_name/Dockerfile"
    local build_context="."

    echo -e "${GREEN}Building $image_name${NC}"
    nerdctl --namespace "$namespace" build -t "mekstrike-$image_name" -f "$dockerfile_path" "$build_context"
    echo -e "${GREEN}Removing $image_name${NC}"
    kubectl delete --ignore-not-found=true -f "k8s/$image_name.yaml"
    echo -e "${GREEN}Deploying $image_name${NC}"
    kubectl apply -f "k8s/$image_name.yaml"
}

echo -e "${GREEN}Ensuring prerequisites${NC}"
kubectl apply -f k8s/mekstrike-base/namespace.yaml
kubectl apply -f k8s/mekstrike-base/mekstrike-instrumentation.yaml
kubectl apply -f k8s/mekstrike-base/mekstrike-config.yaml
kubectl apply -f k8s/mekstrike-base/library-store.yaml
kubectl apply -f k8s/mekstrike-base/battlefield-store.yaml

echo -e "${GREEN}Generating protos${NC}"
protoc --proto_path="." --go_out=. --go_opt=paths=source_relative domain/unit/unit.proto
protoc --proto_path="." --go_out=. --go_opt=paths=source_relative domain/battlefield/battlefield.proto
protoc --proto_path="." --go_out=gamemaster/clients/armybuilder armybuilder/src/main/proto/armybuilder.proto

# Build all images if no parameter is given
if [ $# -eq 0 ]; then
    build_image "library"
    build_image "importer"
    build_image "armybuilder"
    build_image "battlefield"
    build_image "unit"
    build_image "gamemaster"
    build_image "mediaproxy"
    build_image "ui"
    exit 0
fi

# Build command for the specified image
image_name="$1"
build_image "$image_name"
