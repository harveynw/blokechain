# Blokechain

A from-scratch zero dependency* implementation of various parts of the Bitcoin protocol, using the original whitepaper and guides online. Goal is to have a fully functioning wallet and node some time in the future™. Purely educational and not to be trusted in prod.


## <b>internal/script</b>

A fully functioning Bitcoin script interpreter. Can execute P2PK, P2PKH, P2MS, P2SH transactions and anything else allowed by the spec (https://en.bitcoin.it/wiki/Script), except for Locktime opcodes which are still TODO.

## <b>internal/cryptography</b>

This implements secp256k1 ECDSA as well as handling signatures, keypairs and hashing. The clever stuff here is really a port of Andrej Karpathy's excellent blog post: [A from-scratch tour of Bitcoin in Python](http://karpathy.github.io/2021/06/21/blockchain/).

Had to include the /x/crypto module* as RIPEMD160 is not in the stdlib.

 
## <b>internal/chain</b>

These are data structures representing blocks, transactions and merkle trees used in the protocol.

## <b>internal/miner</b>

Block mining functionality.

## <b>internal/wallet</b>

Keypair management and serialisation.