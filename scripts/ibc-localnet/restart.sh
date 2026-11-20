#!/bin/bash
# Rebuild axmd (docker, golang 1.21), rebuild the test image and restart both
# localnet chains on their existing data. Run from Git Bash on the host.
set -euo pipefail
export MSYS_NO_PATHCONV=1
cd "$(dirname "$0")/../.."
ROOT=$(pwd -W)
CTX=${CTX:?set CTX to the docker build context folder}

docker run --rm -v "$ROOT:/src" -v "$HOME/go/pkg/mod:/go/pkg/mod" -v axm-gocache:/root/.cache/go-build \
  -w /src -e GOWORK=off -e GOFLAGS=-mod=mod -e GOTOOLCHAIN=local golang:1.21-bookworm \
  go build -o /src/build/axmd ./cmd/axmd

cp build/axmd "$CTX/"
(cd "$CTX" && docker build -q -t axmd-ibc:dev . >/dev/null)

for c in axm-a axm-b; do
  docker rm -f $c >/dev/null 2>&1 || true
  docker run -d --name $c --network axm-ibc -v $c-data:/data axmd-ibc:dev \
    axmd start --home /data/$c --minimum-gas-prices 0uaxm \
    --grpc.address 0.0.0.0:9090 --rpc.laddr tcp://0.0.0.0:26657 >/dev/null
done
echo "restarted axm-a, axm-b"
