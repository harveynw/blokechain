package chain

import (
	"github.com/harveynw/blokechain/internal/data"
)

// MerkleRoot returns root of merkle tree, duplicating last tx if odd number
func MerkleRoot(txs []Transaction) []byte {
	hashes := make([][]byte, len(txs), len(txs))
	for i, tx := range txs {
		hashes[i] = data.DoubleHash(tx.Encode(-1), false)
	}

	return computeMerkleRoot(hashes)
}

func computeMerkleRoot(hashes [][]byte) []byte {
	// Duplicate lonely hash
	if len(hashes) == 1 {
		hashes = append(hashes, hashes[0])
	}

	if len(hashes) == 2 {
		// Leaf nodes
		return data.DoubleHash(append(hashes[0], hashes[1]...), false)
	}
	// Node
	left, right := computeMerkleRoot(hashes[0:len(hashes)/2]), computeMerkleRoot(hashes[len(hashes)/2:])
	return data.DoubleHash(append(left, right...), false)
}