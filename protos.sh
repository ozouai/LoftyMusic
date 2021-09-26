#!/bin/bash
mkdir -p audiosource/audiosourcepb
protoc -I protos/ --go_out=audiosource/audiosourcepb --go_opt=paths=source_relative --go-grpc_out=audiosource/audiosourcepb --go-grpc_opt=paths=source_relative protos/source.proto