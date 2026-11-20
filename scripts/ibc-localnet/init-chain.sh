#!/bin/sh
# Initialise a single-validator Axiome localnet for IBC tests.
# Usage: init-chain.sh <chain-id>   (home: /data/<chain-id>)
# Keys are the public test mnemonics from scripts/init-testchain.sh.
set -eu

CHAIN=$1
H=/data/$CHAIN
AXMD="axmd --home $H"

# alice: validator + relayer funds (axm1v6xdm93m5s9lvu0ux4k76s2p7hgkzj00u8w5z7)
ALICE="juice dog over thing anger search film document sight fork enrich jungle vacuum grab more sunset winner diesel flock smooth route impulse cheap toward"
# ann: plain user (axm10clqnnaplc06jv4n8a5kyx9qaatuh2au3hrasg)
ANN="output arrange offer advance egg point office silent diamond fame heart hotel rocket sheriff resemble couple race crouch kit laptop document grape drastic lumber"

rm -rf "$H"
$AXMD config set client chain-id "$CHAIN"
$AXMD config set client keyring-backend test
echo "$ALICE" | $AXMD keys add alice --recover --keyring-backend test >/dev/null
echo "$ANN" | $AXMD keys add ann --recover --keyring-backend test >/dev/null

$AXMD init "$CHAIN" --chain-id "$CHAIN" --default-denom uaxm >/dev/null 2>&1
$AXMD genesis add-genesis-account alice 100000000000uaxm --keyring-backend test
$AXMD genesis add-genesis-account ann 100000000uaxm --keyring-backend test

G=$H/config/genesis.json
jq '.app_state.vote.government = ["axm1v6xdm93m5s9lvu0ux4k76s2p7hgkzj00u8w5z7"]' "$G" > "$G.tmp" && mv "$G.tmp" "$G"

$AXMD genesis gentx alice 23529411764uaxm --chain-id "$CHAIN" --keyring-backend test >/dev/null 2>&1
$AXMD genesis collect-gentxs >/dev/null 2>&1
jq '.app_state.referral = input' "$G" /fixtures/referrals.json > "$G.tmp" && mv "$G.tmp" "$G"

# fast blocks, open ports inside the docker network, zero gas price
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/' $H/config/config.toml
sed -i 's#laddr = "tcp://127.0.0.1:26657"#laddr = "tcp://0.0.0.0:26657"#' $H/config/config.toml
sed -i 's/minimum-gas-prices = ".*"/minimum-gas-prices = "0uaxm"/' $H/config/app.toml
sed -i 's/localhost:9090/0.0.0.0:9090/' $H/config/app.toml
$AXMD config set app api.enable true >/dev/null 2>&1 || true

$AXMD genesis validate-genesis
echo "init ok: $CHAIN"
