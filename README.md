# Blokechain

A ground up implementation of Bitcoin in Go, using the original whitepaper and some tutorials online. It includes Script (P2PKH), secp256k1 ECDSA and functional code for running a node + wallet IN DEV.

Transactions on blokechain would be valid on the Bitcoin network, apart from some subtleties in integer encoding, SIGHASH and the use of SHA-256 in place of RIPEMD-160.

# TODO
crypto.go YFromX can simply be refactored to PublicKeys from X or similar, as that is the end goal