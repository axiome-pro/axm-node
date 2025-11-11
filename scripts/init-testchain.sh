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
#       /  \
#    POLLY REF01..15
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

# BOB referrals
# axm196exgjwjs6et5c3fgmjgfrfus3sgxsaw6fe7n2
REF01="organ torch gentle van hollow jealous wheat glide glow stable adult only cattle salute word question interest hat home pink alert deer valve kick"
# axm1j5z4n75yrldsd0fv5wmkcyg9lefnjppru6lk8d
REF02="input cradle science capable kick unfair multiply vague response goddess wish draft remember body liquid consider pill castle fiction build toddler multiply cry favorite"
# axm155ahr0ua8dmw957ln344pplzmpsznuf8xg92jm
REF03="right butter battle repair rocket obey today mimic silent coffee pulp ride assault denial guide strike inquiry floor arctic enhance eye jelly midnight always"
# axm1a04lxj8tpxrykejp8m2a8z8m23tquuc5t4d59c
REF04="chase attitude flower liquid track used visual lawsuit patient document nature throw engage hill auction tired monkey review trouble spoon pudding secret since cheese"
# axm1l9mplvrjtf32zxcmr2tmlyfzgjen2mrp07gsd3
REF05="boy prevent noise summer jeans pumpkin box marine above blush boost umbrella clap guide people yard casino move genre drift orphan provide bread walnut"
# axm1l3ceys5hu4e54px4d8heust49aww8cj2peyja5
REF06="amused orient guide giant train visit gas hospital object outdoor vacuum fantasy arena quantum thought two morning mango install neither sweet lady palace script"
# axm17xalfwp7nqsz6rep6lk7yhu4u92a89xgl8hdma
REF07="sure protect siege junior engine photo chimney select dad walnut solve permit good bus space can sniff toss seven omit weekend swamp vanish anger"
# axm13w9x9zq0t59v5erhfnhqljevwmcl32a50sul9d
REF08="cool eagle abstract narrow picnic issue report fat loop share rhythm negative talk track unfair ketchup design loyal wedding please easily rigid basic include"
# axm1g0ayvgakp9qhyqqz0u2s0jm9076g87racw95q7
REF09="soldier make wood scrap lottery fiber income dish enrich jar pull treat soda tide crouch fade cool chaos unable dizzy country toddler mammal text"
# axm1naragks4qcwapc7wqp464kxt5xc2crz9k6qgkx
REF10="toy crop forum warm topic critic chair gauge neither pave vast ill march fatigue attend pumpkin wide gorilla wrap shrug protect output napkin seek"
# axm1h696aa4x4eyh9gjy6n5dwvcsj3wurpjgx2356a
REF11="monkey mask gorilla series marriage dentist illness scrap party orbit limb unusual rug traffic include crouch artist insect stage cargo disagree author label dawn"
# axm1cuv689ptlfgkehya68fekf58lk67f07g5wx4n8
REF12="prize erase obey fork depart area soul front nasty claw lake whip topic decade hawk enhance mixed dumb glare system chase wise such west"
# axm1anwnf3usfy3hqjm7tu5ugjrr2q3p43w7c8747w
REF13="ritual alley fish gentle silly short laptop ocean crouch change emotion wool unknown bubble soda ride agree hello cry dice luggage fatigue install town"
# axm1jc3y5zhwzj2ejfvn2277770hnp2advm9jchh4n
REF14="antenna example moon foam finish scrub indoor three firm crack moon split ivory trash time antique drink ostrich flame harbor float add two gadget"
# axm10e9vzc7zcn86ut55crf4erlm3fltj8e45u3yxc
REF15="nerve mom brush enforce senior suit animal picture normal chat slide curve nurse sweet rely prevent kick draft steak size odor witness young town"


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


yes "$REF01" | $axmd_BIN keys add REF01 --recover
yes "$REF02" | $axmd_BIN keys add REF02 --recover
yes "$REF03" | $axmd_BIN keys add REF03 --recover
yes "$REF04" | $axmd_BIN keys add REF04 --recover
yes "$REF05" | $axmd_BIN keys add REF05 --recover
yes "$REF06" | $axmd_BIN keys add REF06 --recover
yes "$REF07" | $axmd_BIN keys add REF07 --recover
yes "$REF08" | $axmd_BIN keys add REF08 --recover
yes "$REF09" | $axmd_BIN keys add REF09 --recover
yes "$REF10" | $axmd_BIN keys add REF10 --recover
yes "$REF11" | $axmd_BIN keys add REF11 --recover
yes "$REF12" | $axmd_BIN keys add REF12 --recover
yes "$REF13" | $axmd_BIN keys add REF13 --recover
yes "$REF14" | $axmd_BIN keys add REF14 --recover
yes "$REF15" | $axmd_BIN keys add REF15 --recover

$axmd_BIN init test --chain-id demo --default-denom uaxm

# update genesis
$axmd_BIN genesis add-genesis-account alice 100000000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account ann 100000000uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account ben 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account bob 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account den 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account miranda 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account polly 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account tom 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account joe 100uaxm --keyring-backend test

$axmd_BIN genesis add-genesis-account REF01 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF02 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF03 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF04 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF05 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF06 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF07 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF08 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF09 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF10 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF11 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF12 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF13 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF14 100uaxm --keyring-backend test
$axmd_BIN genesis add-genesis-account REF15 100uaxm --keyring-backend test

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
