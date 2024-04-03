#!/usr/bin/env bash

axmd_BIN=$(which axmd)
validator=$(axmd q staking validators --output json | jq ".validators[0].operator_address" | sed "s/\"//" | sed "s/\"//")

echo $validator

$axmd_BIN tx staking delegate $validator 300000000uaxm --gas auto --gas-adjustment=1.4 --yes --from ben
$axmd_BIN tx staking delegate $validator 8500000000uaxm --gas auto --gas-adjustment=1.4 --yes --from bob
$axmd_BIN tx staking delegate $validator 8500000000uaxm --gas auto --gas-adjustment=1.4 --yes --from den
$axmd_BIN tx staking delegate $validator 8500000000uaxm --gas auto --gas-adjustment=1.4 --yes --from miranda
#
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF01
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF02
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF03
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF04
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF05
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF06
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF07
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF08
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF09
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF10
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF11
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF12
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF13
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF14
$axmd_BIN tx staking delegate $validator 850000000uaxm --gas auto --gas-adjustment=1.4 --yes --from REF15
