#!/usr/bin/env bash

set -e

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    # this regex checks if a proto file has its go_package set to github.com/cosmosregistry/example/api/...
    # gogo proto files SHOULD ONLY be generated if this is false
    # we don't want gogo proto to run for proto files which are natively built for google.golang.org/protobuf
    if grep -q "option go_package" "$file" && grep -H -o -c 'option github.com/axiome-pro/axm-node' "$file" | grep -q ':0$'; then
      echo "Generating buf for " $file
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

echo "Generating pulsar proto code"
buf generate --template buf.gen.pulsar.yaml

cd ..

cp -r github.com/axiome-pro/axm-node/* ./
rm -rf api && mkdir api
mv axiome ./api

rm -rf github.com
sed -i 's/github.com\/cometbft\/cometbft\/proto\/tendermint\/types/cosmossdk.io\/api\/tendermint\/types/' api/axiome/staking/v1beta1/staking.pulsar.go
sed -i 's/github.com\/cometbft\/cometbft\/abci\/types/cosmossdk.io\/api\/tendermint\/abci/' api/axiome/staking/v1beta1/staking.pulsar.go

rm -r github.com

## Apply CosmWasm replacements to all pulsar.go files
#find api -name "*.pulsar.go" -exec sed -i 's/github.com\/CosmWasm\/wasmd\/x\/wasm\/types/github.com\/CosmWasm\/wasmd\/x\/wasm\/types/g' {} \;
#
## Specifically handle module.proto imports for CosmWasm
#find api -name "*.pulsar.go" -exec sed -i 's/github.com\/CosmWasm\/wasmd\/x\/wasm\/types\/module/github.com\/CosmWasm\/wasmd\/x\/wasm\/types/g' {} \;

#sed -i 's/github.com\/cometbft\/cometbft\/proto\/tendermint\/types/cosmossdk.io\/api\/tendermint\/types/' api/delegating/v1beta1/staking.pulsar.go
#sed -i 's/github.com\/cometbft\/cometbft\/abci\/types/cosmossdk.io\/api\/tendermint\/abci/' api/delegating/v1beta1/de.pulsar.go
