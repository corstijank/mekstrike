#!/bin/sh
protoc --proto_path="." --go_out=. --go_opt=paths=source_relative common/go/unit/stats.proto
protoc --proto_path="." --go_out=gamemaster/clients/armybuilder armybuilder/src/main/proto/armybuilder.proto
