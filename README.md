# Axiome Chain
This repository contains an Axiome chain node's source code.
Axiome is an ecosystem that has been built on its own L1 solution, Axiome Chain,
and its sole token, AXM. New tokens are mined through time-tested DPoS
delegation of AXM. The delegation reward rate is floating based on
bonded / free tokens ratio.

`axmd` uses the **0.50.3** version of [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk).

## How to Use
### Installation (for test purposes)
Install and run demo `axmd` node:

```sh
git clone git@github.com:axiome-pro/axm-node.git
cd axm-node
make install # install the axmd binary
make init # initialize the demo chain
axmd start # start the demo chain
```
You can find test accounts in `scripts/init-testchain.sh`

### Installation (for production use)

Please refer to [this article](https://docs.axiomeinfo.org/how-to/start-blockchain-node-on-ubuntu) for production node installation instructions

## Useful links
* [Axiome Project Documentation](https://axiomeinfo.org/)
* [Axiome Chain Documentation](https://docs.axiomeinfo.org/)
* [Cosmos-SDK Documentation](https://docs.cosmos.network/)
