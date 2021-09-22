package node 

import (
	"fmt"
	"net/rpc/jsonrpc"
)

// Node is our server object for holding our version of the blockchain and communicating with other nodes
type Node struct {
	LocalChain *Blockchain
}

