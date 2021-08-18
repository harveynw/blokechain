package chain

import (
	"fmt"
	"github.com/harveynw/blokechain/internal/data"
)

type Block struct {
	header BlockHeader
	txs []Transaction
}

type BlockHeader struct {
	prevBlockHash []byte
	merkleRoot []byte
	timestamp uint32
	difficultyTarget []byte
	nonce []byte
}

func (bh BlockHeader) Encode() []byte {
	enc := make([]byte, 0)

	// Version
	enc = append(enc, 0x02, 0x00, 0x00, 0x01)
	
	// Previous Block Hash (32 Bytes)
	enc = append(enc, bh.prevBlockHash[0:32]...)

	// Merkle Root (32 Bytes)
	enc = append(enc, bh.merkleRoot[0:32]...)

	// Timestamp (seconds from Unix Epoch, 4 Bytes)
	enc = append(enc, data.EncodeInt(bh.timestamp, 4)...)

	// Difficulty Target
	enc = append(enc, data.EncodeInt(int(bh.difficultyTarget), 4)...)

}