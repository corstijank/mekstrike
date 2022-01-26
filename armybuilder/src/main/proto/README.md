# Common Protobuf Definitions for Mekstrike

## ArmyBuilder

```bash
# Generate Go package for gamemaster; place it in the abclient package
protoc --proto_path="../../../.." --go_out=../../../../gamemaster/abclient armybuilder.proto
```
