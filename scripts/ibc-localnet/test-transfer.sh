#!/bin/sh
# IBC round trip between axm-a and axm-b over transfer/channel-0.
# Run inside the axm-a container:  docker exec -i axm-a sh < test-transfer.sh
set -eu
ANN=axm10clqnnaplc06jv4n8a5kyx9qaatuh2au3hrasg
A="--node tcp://axm-a:26657"
B="--node tcp://axm-b:26657"
KR="--home /data/axm-a --keyring-backend test"

bal() { axmd q bank balances "$2" $1 -o json | jq -c '[.balances[] | {(.denom): .amount}] | add // {}'; }
tx() {  # chain-flags chain-id args...
  n=$1; c=$2; shift 2
  out=$(axmd tx ibc-transfer transfer transfer channel-0 "$@" --from ann $KR $n --chain-id $c \
        --gas 300000 --fees 0uaxm -y -o json)
  echo "$out" | jq -r '"  tx \(.txhash) code=\(.code)"'
  sleep 3
  h=$(echo "$out" | jq -r .txhash)
  axmd q tx $h $n -o json | jq -r '"  в блоке \(.height), code=\(.code)"'
}

# ann's key lives in axm-a's keyring; the same mnemonic owns the same address on axm-b
echo "== до"
echo "  A: $(bal "$A" $ANN)"; echo "  B: $(bal "$B" $ANN)"

echo "== A → B: 1 AXM"
tx "$A" axm-a "$ANN" 1000000uaxm
sleep 15
echo "  A: $(bal "$A" $ANN)"; echo "  B: $(bal "$B" $ANN)"

IBCD=$(axmd q bank balances $ANN $B -o json | jq -r '.balances[] | select(.denom|startswith("ibc/")) | .denom')
echo "  единица на B: $IBCD"
axmd q ibc-transfer denom-trace "${IBCD#ibc/}" $B -o json | jq -c .

ESC=$(axmd q ibc-transfer escrow-address transfer channel-0 $A)
echo "  escrow на A ($ESC): $(bal "$A" $ESC)"

echo "== B → A: 0.4 AXM обратно"
tx "$B" axm-b "$ANN" "400000$IBCD"
sleep 15
echo "  A: $(bal "$A" $ANN)"; echo "  B: $(bal "$B" $ANN)"
echo "  escrow на A: $(bal "$A" $ESC)"
echo "  всего ibc-единицы на B: $(axmd q bank total-supply-of "$IBCD" $B -o json | jq -c .)"
