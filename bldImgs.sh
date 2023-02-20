#!/bin/bash
nerdctl --namespace k8s.io build -t mekstrike-library -f library/Dockerfile .
nerdctl --namespace k8s.io build -t mekstrike-importer -f library/importer/Dockerfile .
nerdctl --namespace k8s.io build -t mekstrike-mediaproxy -f mediaproxy/Dockerfile .
nerdctl --namespace k8s.io build -t mekstrike-gamemaster -f gamemaster/Dockerfile .
nerdctl --namespace k8s.io build -t mekstrike-armybuilder -f armybuilder/Dockerfile .
nerdctl --namespace k8s.io build -t mekstrike-battlefield -f battlefield/Dockerfile battlefield
nerdctl --namespace k8s.io build -t mekstrike-unit -f unit/Dockerfile .
nerdctl --namespace k8s.io build -t mekstrike-ui -f ui/Dockerfile ui/
