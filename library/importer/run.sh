#!/bin/bash
dapr run --app-id library \
--components-path ../components \
go run importer.go