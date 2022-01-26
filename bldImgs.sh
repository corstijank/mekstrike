#!/bin/bash
docker build -t mekstrike-library -f library/Dockerfile .
docker build -t mekstrike-importer -f library/importer/Dockerfile .
docker build -t mekstrike-gamemaster -f gamemaster/Dockerfile .
docker build -t mekstrike-armybuilder -f armybuilder/Dockerfile .
docker build -t mekstrike-battlefield -f battlefield/Dockerfile battlefield
docker build -t mekstrike-unit -f unit/Dockerfile .
docker build -t mekstrike-ui -f ui/Dockerfile ui/
