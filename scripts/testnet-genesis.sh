#!/usr/bin/env bash

rm -r ~/.axmd || true
axmd_BIN=$(which axmd)

NEWLINE=$'\n'
DIR="$( dirname -- "${BASH_SOURCE[0]}"; )";
DIR="$( realpath -e -- "$DIR"; )";

# ACCS
AXM_NODE_1="fabric soup similar kitten surround purchase forget gesture salad humor pencil arch wait kingdom pride kite ridge trouble cat practice reject medal increase insect"

# configure axmd
$axmd_BIN config set client chain-id testnet-12
$axmd_BIN config set client keyring-backend test

yes "$AXM_NODE_1" | $axmd_BIN keys add axm-node-1 --recover

$axmd_BIN init axm-node-1 --chain-id testnet-12 --default-denom uaxm
# gen axm acc
# update genesis
$axmd_BIN genesis add-genesis-account axm10648lnjrzvfpzng3ryk4jt5y84jd5agh6yjrvq 5000000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account axm19r09spfv5qszxc53hyqgmd7hugam70c5lel9uz 5882352941177uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account axm1ytal6q2qc3rkha3x24m4ynh9nvzmpkm8lucq00 5858823529412uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account axm1a3elvwqv9339heznngtm3sawzh8ertrg4xfyax 23529411765uaxm --keyring-backend test
# create default validator
$axmd_BIN genesis gentx axm-node-1 23529411765uaxm --chain-id testnet-12
$axmd_BIN genesis collect-gentxs

SOURCE_GENESIS=~/.axmd/config/genesis.json
TMP_GENESIS1=~/.axmd/config/tmp_genesis1.json
TMP_GENESIS2=~/.axmd/config/tmp_genesis2.json
RESULT_GENESIS=~/.axmd/config/tmp_genesis.json

MEGASTATUS_TIME='2024-12-31T23:59:59.9999999999Z'
DOWNGRADES='{"account": "axm1ytal6q2qc3rkha3x24m4ynh9nvzmpkm8lucq00","current": "STATUS_MEGA","time": "'$MEGASTATUS_TIME'"}'
DOWNGRADES=$DOWNGRADES',{"account": "axm19qkyrg3qtscnlzlpsuhsex2qumlyxayxga2jj2","current": "STATUS_MEGA","time": "'$MEGASTATUS_TIME'"}'
DOWNGRADES=$DOWNGRADES',{"account": "axm1fjl67wz3vydwettad5fj9ntguxj7tjxmd9pr53","current": "STATUS_MEGA","time": "'$MEGASTATUS_TIME'"}'
DOWNGRADES=$DOWNGRADES',{"account": "axm1wj3jhzs7n5qlaq7tsy6uw2y454u50hwmctkvt6","current": "STATUS_MEGA","time": "'$MEGASTATUS_TIME'"}'

TOP_LEVEL='{"address": "axm1ytal6q2qc3rkha3x24m4ynh9nvzmpkm8lucq00", "status": "STATUS_NEW"}'
TOP_LEVEL=$TOP_LEVEL',{"address": "axm19r09spfv5qszxc53hyqgmd7hugam70c5lel9uz", "status": "STATUS_NEW"}'

jq ".app_state.referral.top_level_accounts = [$TOP_LEVEL]" $SOURCE_GENESIS > $TMP_GENESIS1
jq '.app_state.referral.other_accounts = input.other_accounts' $TMP_GENESIS1 "$DIR/referrals_testnet.json" > $TMP_GENESIS2
jq '.app_state.distribution.params.community_tax = "0.000000000000000000"' $TMP_GENESIS2 > $TMP_GENESIS1
jq '.app_state.distribution.params.validator_commission_rate = "0.050000000000000000"' $TMP_GENESIS1 > $TMP_GENESIS2
jq '.app_state.referral.params.status_downgrade_period = 14400' $TMP_GENESIS2 > $TMP_GENESIS1
jq '.app_state.staking.params.unbonding_time = "3600s"' $TMP_GENESIS1 > $TMP_GENESIS2
jq ".app_state.referral.downgrades = [$DOWNGRADES]" $TMP_GENESIS2 > $TMP_GENESIS1

cp $TMP_GENESIS1 $RESULT_GENESIS