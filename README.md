# Blokechain

A ground up implementation of Bitcoin in Go, using the original whitepaper and some tutorials online. It includes Script (P2PKH), secp256k1 ECDSA and functional code for running a node + wallet IN DEV.

Transactions on blokechain would be valid on the Bitcoin network, apart from some subtleties in integer encoding, SIGHASH and the use of SHA-256 in place of RIPEMD-160.

# TODO
- Wallet command line access
- Node server logic
- Miner logic
- Implementing all of bitcoin script
- OP_DAVE