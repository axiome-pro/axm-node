#!/usr/bin/env bash

axmd_BIN=$(which axmd)
validator=$(axmd q staking validators --output json | jq ".validators[0].operator_address" | sed "s/\"//" | sed "s/\"//")

echo $validator

$axmd_BIN tx staking delegate $validator 300000000uaxm --gas auto --gas-adjustment=1.4 --yes --from ben
$axmd_BIN tx staking delegate $validator 8500000000uaxm --gas auto --gas-adjustment=1.4 --yes --from bob
$axmd_BIN tx staking delegate $validator 8500000000uaxm --gas auto --gas-adjustment=1.4 --yes --from den
$axmd_BIN tx staking delegate $validator 8500000000uaxm --gas auto --gas-adjustment=1.4 --yes --from miranda
#
## configure axmd
#$axmd_BIN config set client chain-id demo
#$axmd_BIN config set client keyring-backend test
#$axmd_BIN keys add alice
#$axmd_BIN keys add bob
#$axmd_BIN init test --chain-id demo --default-denom uaxm
## update genesis
#$axmd_BIN genesis add-genesis-account alice 10000000uaxm --keyring-backend test
#$axmd_BIN genesis add-genesis-account bob 1000uaxm --keyring-backend test
## create default validator
#$axmd_BIN genesis gentx alice 1000000uaxm --chain-id demo
#$axmd_BIN genesis collect-gentxs
