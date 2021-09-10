package chain

import (
	"github.com/harveynw/blokechain/internal/data"
)

// Block data structure for forming blockchain
type Block struct {
	header *BlockHeader
	txs []Transaction
}

// BlockHeader data structure for summarising contents of a block
type BlockHeader struct {
	prevBlockHash []byte
	merkleRoot []byte
	timestamp uint32
	difficultyTarget Difficulty
	nonce uint32
}

// Encode serializes the block header, used for mining
func (bh BlockHeader) Encode() []byte {
	enc := make([]byte, 0)

	// Version
	enc = append(enc, 0x02, 0x00, 0x00, 0x01)
	
	// Previous Block Hash (32 Bytes)
	enc = append(enc, bh.prevBlockHash[0:32]...)

	// Merkle Root (32 Bytes)
	enc = append(enc, bh.merkleRoot[0:32]...)

	// Timestamp (seconds from Unix Epoch, 4 Bytes)
	enc = append(enc, data.EncodeInt(int(bh.timestamp), 4)...)

	// Difficulty Target
	enc = append(enc, bh.difficultyTarget.Encode()...)

	// Nonce
	enc = append(enc, data.EncodeInt(int(bh.nonce), 4)...)

	return enc
}

