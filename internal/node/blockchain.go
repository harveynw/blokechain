package node

import (
	"fmt"
	"github.com/harveynw/blokechain/internal/chain"
)

// Blockchain is a directed graph structure for storing competing chains
type Blockchain struct {
	Heights map[string]int // Stores all block heights
	Nodes map[string]chain.Block // Stores all blocks
	Children map[string]map[string]bool // block hash -> children block hashes
	Parent map[string]string // block hash -> parent block hash
}

// AddBlock attempts to add a block to the chain, provided the parent is stored
func (chain *Blockchain) AddBlock(block chain.Block) (success bool, parentRequired string) {
	current := hex(block.Header.BlockHash())
	parent := hex(block.Header.PrevBlockHash)

	if _, ok := chain.Nodes[parent]; ok {
		// Add block to nodes
		chain.Nodes[current] = block

		// Attaching to parent
		chain.Children[parent][current] = true 
		chain.Parent[current] = parent

		// Store height (useful for otherwise NP-hard longest chain)
		chain.Heights[current] = chain.computeHeight(current)

		return true, ""
	}

	// Parent not in chain, return false
	return false, parent
}

func (chain *Blockchain) computeHeight(block string) int {
	if _, ok := chain.Nodes[block]; !ok {
		// Not in chain
		return -1
	}

	height := 0
	for {
		if parent, hasParent := chain.Parent[block]; !hasParent {
			// Reached genesis
			break
		} else {
			block = parent
			height++
		}
	}

	return height
}

func hex(b []byte) string {
	return fmt.Sprintf("%x", b)
}

