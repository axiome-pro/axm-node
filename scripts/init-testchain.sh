#!/usr/bin/env bash

rm -r ~/.axmd || true
axmd_BIN=$(which axmd)
NEWLINE=$'\n'
DIR="$( dirname -- "${BASH_SOURCE[0]}"; )";
DIR="$( realpath -e -- "$DIR"; )";

#
# TEST referral structure
#
#     ALICE
#     /   \
#    ANN  BEN ----
#         /  \    \
#       BOB  DEN  MIRANDA
#       /
#    POLLY
#     /
#   TOM
#

# validator: axmvaloper1v6xdm93m5s9lvu0ux4k76s2p7hgkzj00kw2z6h
# axm1v6xdm93m5s9lvu0ux4k76s2p7hgkzj00u8w5z7
ALICE="juice dog over thing anger search film document sight fork enrich jungle vacuum grab more sunset winner diesel flock smooth route impulse cheap toward"
# axm10clqnnaplc06jv4n8a5kyx9qaatuh2au3hrasg
ANN="output arrange offer advance egg point office silent diamond fame heart hotel rocket sheriff resemble couple race crouch kit laptop document grape drastic lumber"
# axm1rkk2x2jedsdkmckyrxzjx9tdtavnqwlpqf3yts
BEN="keep liar demand upon shed essence tip undo eagle run people strong sense another salute double peasant egg royal hair report winner student diamond"
# axm157lyytcwlk4lg2h7tvssapf9x45g7re5mlv96r
BOB="gesture inject test cycle original hollow east ridge hen combine junk child bacon zero hope comfort vacuum milk pitch cage oppose unhappy lunar seat"
# axm15ltmumw44fljaj792z6ymzyt54uaqz6c54u8wk
DEN="copper push brief egg scan entry inform record adjust fossil boss egg comic alien upon aspect dry avoid interest fury window hint race symptom"
# axm1042mz3m4c6cxeucjd609e8eajqua8z0crlk3td
MIRANDA="pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear"
# axm1f3gx40xsh34u35xfyymedkt55mlfkxqzrd6gfu
POLLY="earn front swamp dune level clip shell aware apple spare faith upset flip local regret loud suspect view heavy raccoon satisfy cupboard harbor basic"
# axm14nnlxrtqchasxa7u6u0vg7pmqpu8jjxg74x4u3
TOM="maximum display century economy unlock van census kite error heart snow filter midnight usage egg venture cash kick motor survey drastic edge muffin visual"
# axm1lgh4mzy5es9qs4zqhwfrz03sjuuwydttr5xjap
JOE="apart acid night more advance december weather expect pause taxi reunion eternal crater crew lady chaos visual dynamic friend match glow flash couple tumble"

# configure axmd
$axmd_BIN config set client chain-id demo
$axmd_BIN config set client keyring-backend test

yes "$ALICE" | $axmd_BIN keys add alice --recover
yes "$ANN" | $axmd_BIN keys add ann --recover
yes "$BEN" | $axmd_BIN keys add ben --recover
yes "$BOB" | $axmd_BIN keys add bob --recover
yes "$DEN" | $axmd_BIN keys add den --recover
yes "$MIRANDA" | $axmd_BIN keys add miranda --recover
yes "$POLLY" | $axmd_BIN keys add polly --recover
yes "$TOM" | $axmd_BIN keys add tom --recover
yes "$JOE" | $axmd_BIN keys add joe --recover

$axmd_BIN init test --chain-id demo --default-denom uaxm

# update genesis
$axmd_BIN genesis add-genesis-account alice 100000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account ann 10000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account ben 10000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account bob 10000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account den 10000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account miranda 10000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account polly 10000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account tom 10000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account joe 1000000000uaxm --keyring-backend test

# vote gov module params
jq '.app_state.vote.government = ["axm1v6xdm93m5s9lvu0ux4k76s2p7hgkzj00u8w5z7"]' ~/.axmd/config/genesis.json > ~/.axmd/config/tmp_genesis.json
mv ~/.axmd/config/tmp_genesis.json ~/.axmd/config/genesis.json -f


# create default validator
$axmd_BIN genesis gentx alice 23529411764uaxm --chain-id demo

$axmd_BIN genesis collect-gentxs

jq '.app_state.referral = input' ~/.axmd/config/genesis.json "$DIR/../tests/fixtures/referrals.json" > ~/.axmd/config/tmp_genesis.json
mv ~/.axmd/config/tmp_genesis.json ~/.axmd/config/genesis.json -f

$axmd_BIN config set app api.enable true
$axmd_BIN config set app api.swagger true

sed -i 's/127.0.0.1/0.0.0.0/' ~/.axmd/config/config.toml
sed -i 's/localhost/0.0.0.0/' ~/.axmd/config/app.toml
