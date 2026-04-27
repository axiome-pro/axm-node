#!/usr/bin/env bash

cd /
apt-get update && apt-get install -y git make gcc
git config --global --add safe.directory /axm-node
cd /axm-node
rm -rf builds

for os in linux
do
  CGO_ENABLED=1 GOARCH=amd64 GOOS=$os OUTPUT_DIR=builds/${GOOS}_${GOARCH} GOFLAGS="-trimpath -o=${OUTPUT_DIR}/" /bin/sh -c 'mkdir -p ${OUTPUT_DIR} && make build'
done
