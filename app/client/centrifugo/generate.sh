#!/bin/bash

set -e

rm -rf /tmp/proto/*
mkdir -p /tmp/proto/apiproto

wget -O /tmp/proto/api.template.proto https://raw.githubusercontent.com/centrifugal/centrifuge/master/misc/proto/api.template.proto
wget -O /tmp/proto/client.template.proto https://raw.githubusercontent.com/centrifugal/centrifuge/master/misc/proto/client.template.proto

gomplate -f /tmp/proto/api.template.proto > /tmp/proto/api.proto
gomplate -f /tmp/proto/client.template.proto > /tmp/proto/client.proto

rm /tmp/proto/api.template.proto
rm /tmp/proto/client.template.proto

cd /tmp/proto/apiproto && protoc --proto_path=/tmp/proto  --gogofaster_out=plugins=grpc:. /tmp/proto/api.proto
cd /tmp/proto && protoc --proto_path=/tmp/proto  --gogofaster_out=plugins=grpc:. /tmp/proto/client.proto
