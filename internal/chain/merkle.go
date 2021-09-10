package chain

import (
	"github.com/harveynw/blokechain/internal/data"
)

// MerkleRoot returns root of merkle tree, duplicating last tx if odd number
func (b *Block) MerkleRoot() []byte {
	hashes := make([][]byte, 0, len(b.txs))
	for i, tx := range b.txs {
		hashes[i] = data.DoubleHash(tx.Encode(-1), false)
	}

	return computeMerkleRoot(hashes)
}

func computeMerkleRoot(hashes [][]byte) []byte {
	nodes := len(hashes)

	// Duplicate lonely hash
	if nodes == 1 {
		hashes = append(hashes, hashes[0])
	}

	if nodes == 2 {
		// Leaf nodes
		return data.DoubleHash(append(hashes[0], hashes[1]...), false)
	}
	// Node
	left, right := computeMerkleRoot(hashes[0:nodes/2]), computeMerkleRoot(hashes[nodes/2:nodes])
	return data.DoubleHash(append(left, right...), false)
}