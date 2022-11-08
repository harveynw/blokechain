package chain

import (
	"fmt"
	"math/big"
	"time"
	"github.com/harveynw/blokechain/internal/cryptography"
)

// Block data structure for forming blockchain
type Block struct {
	Header *BlockHeader
	txs []Transaction
}

// BlockHeader data structure for summarising contents of a block
type BlockHeader struct {
	PrevBlockHash []byte
	merkleRoot []byte
	timestamp uint32
	DifficultyTarget Difficulty
	Nonce uint32
}

// Encode serializes the block header, used for mining
func (bh BlockHeader) Encode() []byte {
	enc := make([]byte, 0)

	// Version
	enc = append(enc, 0x02, 0x00, 0x00, 0x01)
	
	// Previous Block Hash (32 Bytes)
	enc = append(enc, bh.PrevBlockHash[0:32]...)

	// Merkle Root (32 Bytes)
	enc = append(enc, bh.merkleRoot[0:32]...)

	// Timestamp (seconds from Unix Epoch, 4 Bytes)
	enc = append(enc, cryptography.EncodeInt(int(bh.timestamp), 4)...)

	// Difficulty Target
	enc = append(enc, bh.DifficultyTarget.Encode()...)

	// Nonce
	enc = append(enc, cryptography.EncodeInt(int(bh.Nonce), 4)...)

	return enc
}

// BlockHash hashes the block header
func (bh BlockHeader) BlockHash() []byte {
	enc := bh.Encode()
	return cryptography.DoubleHash(enc, false)
}

// IncrementNonce increases the nonce field of the block to ensure it's hash changes
func (bh BlockHeader) IncrementNonce() BlockHeader {
	bh.Nonce++

	// Exhausted, adjust timestamp
	if bh.Nonce == 0 {
		fmt.Println("Nonce exhausted")
		// Y2K 2.0 rip :(
		currentTimestamp := uint32(time.Now().Unix())
		if bh.timestamp == currentTimestamp {
			bh.timestamp++
		} else {
			bh.timestamp = currentTimestamp
		}
	}

	return bh
}

// Genesis returns the genesis block
func Genesis(diff int64) Block {
	coinbaseTx := Transaction{
		version: 1,
		txIn: []TransactionInput{
			{
				prevTransaction: make([]byte, 32),
				prevIndex: 0,
				prevTransactionPubKey: []byte{0x1},
			},
		},
		txOut: []TransactionOutput{
			{
				amount: 5000000001,
				scriptPubKey: []byte{118,169,20,162,200,202,90,247,28,46,166,63,168,79,33,36,17,56,57,203,205,30,195,136,172},
			},
		},
	}

	genesisDifficulty, _ := DecodeDifficulty([]byte{0x1e, 0xff, 0xff, 0xff})
	//fmt.Printf("GEN DIFF %X \n", genesisDifficulty.targetBytes)
	genesisDifficulty = genesisDifficulty.Div(big.NewInt(diff))
	//fmt.Printf("GEN DIFF %X \n", genesisDifficulty.targetBytes)
	//genesisDifficulty, _ := DecodeDifficulty([]byte{0x1e, 0x19, 0xFF, 0xFF})
	// fmt.Printf("%X \n", genesisDifficulty.targetBytes)
	merkleRoot := MerkleRoot([]Transaction{coinbaseTx})
	bh := BlockHeader{
		PrevBlockHash: make([]byte, 32),
		merkleRoot: merkleRoot,
		timestamp: 1631303394,
		DifficultyTarget: genesisDifficulty,
		Nonce: 0,
	}

	//fmt.Printf("Merkle root %x \n", merkleRoot)

	return Block{
		Header: &bh,
		txs: []Transaction{coinbaseTx},
	}
}
