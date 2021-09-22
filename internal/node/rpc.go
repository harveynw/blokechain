package node

// Docstrings from https://developer.bitcoin.org/reference/rpc/

import (
	"errors"
	"github.com/harveynw/blokechain/internal/chain"
)

// GetBestBlockHash Returns the hash of the best (tip) block in the most-work fully-validated chain.
func (node *Node) GetBestBlockHash(empty struct{}, best *string) error {
	maxHeight := -1
	var bestBlock string

	for block, height := range node.LocalChain.Heights {
		if height > maxHeight {
			maxHeight = height
			bestBlock = block
		}
	}

	if maxHeight == -1 {
		return errors.New("Blockchain is empty")
	}

	*best = bestBlock
	return nil
}

// GetBlock returns an Object with information about block ‘hash’ and information about each transaction.
func (node *Node) GetBlock(blockhash string, block *chain.Block) error {
	if found, inChain := node.LocalChain.Nodes[blockhash]; !inChain {
		return errors.New("Block not found")
	} else {
		*block = found 
	}

	return nil
}

// GetBlockCount Returns the height of the most-work fully-validated chain.
func (node *Node) GetBlockCount(empty struct{}, count *int) error {
	maxHeight := -1

	for _, height := range node.LocalChain.Heights {
		if height > maxHeight {
			maxHeight = height
		}
	}

	if maxHeight == -1 {
		return errors.New("Blockchain is empty")
	}

	*count = maxHeight
	return nil
}

// GetBlockHash Returns hash of block in best-block-chain at height provided.
func (node *Node) GetBlockHash(height int, hash *string) error {
	var best *string;
	err := node.GetBestBlockHash(struct{}{}, best)

	if err != nil {
		return err
	}
	if node.LocalChain.Heights[*best] < height || height < 0 {
		return errors.New("No blocks at given height")
	}

	// Walk back until height is met
	current := *best
	for {
		if node.LocalChain.Heights[current] == height {
			*hash = current
			return nil
		}
		current = node.LocalChain.Parent[current]
	}

	return errors.New("Blockchain malformed")
}

// GetBlockHeader returns an Object with information about blockheader ‘hash’.
func (node *Node) GetBlockHeader(hash string, header *chain.BlockHeader) error {
	if found, inChain := node.LocalChain.Nodes[hash]; !inChain {
		return errors.New("Block not found")
	} else {
		*header = *found.Header
	}

	return nil
}

