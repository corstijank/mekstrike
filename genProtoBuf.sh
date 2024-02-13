#!/bin/sh
protoc --proto_path="." --go_out=. --go_opt=paths=source_relative domain/unit/unit.proto
protoc --proto_path="." --go_out=. --go_opt=paths=source_relative domain/battlefield/battlefield.proto

protoc --proto_path="." --go_out=gamemaster/clients/armybuilder armybuilder/src/main/proto/armybuilder.proto