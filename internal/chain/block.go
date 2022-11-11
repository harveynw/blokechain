package chain

import (
	"bytes"
)

// Block data structure for forming blockchain
type Block struct {
	Header *BlockHeader
	txs []Transaction
}

func (block Block) Encode() []byte {
	b := make([]byte, 0)
	b = append(b, []byte{0xD9, 0xB4, 0xBE, 0xF9}...) // Magic no
	b = append(b, []byte{0x00, 0x00, 0x00, 0x00}...) // Blocksize to be amended
	b = append(b, block.Header.Encode()...) // Block header

	n_txs := VarInt{int64(len(block.txs))}
	b = append(b, n_txs.EncodeVarInt()...)

	for _, tx := range(block.txs) {
		b = append(b, tx.Encode(-1)...)
	}

	// Amend Blocksize
	blocksize := EncodeInt(int64(len(b) - 8), 4)
	for i, v := range blocksize {
		b[i+4] = v
	}

	return b
}

func DecodeBlock(b []byte) Block {
	if bytes.Compare(b[0:4], []byte{0xD9, 0xB4, 0xBE, 0xF9}) != 0 {
		panic("Block malformed")
	}

	// Discard magic no, blocksize
	b = b[8:]

	bh, b := DecodeNextBlockHeader(b)

	n_txs, b := DecodeNextVarInt(b)

	txs := make([]Transaction, 0)
	for i := int64(0); i < n_txs.val; i++ {
		var tx Transaction
		tx, b = DecodeNextTransaction(b)
		txs = append(txs, tx)
	}

	return Block{Header: &bh, txs: txs}
}

// // Genesis returns the genesis block
// func Genesis(diff int64) Block {
// 	coinbaseTx := Transaction{
// 		version: 1,
// 		txIn: []TransactionInput{
// 			{
// 				prevTransaction: make([]byte, 32),
// 				prevIndex: 0,
// 				prevTransactionPubKey: []byte{0x1},
// 			},
// 		},
// 		txOut: []TransactionOutput{
// 			{
// 				amount: 5000000001,
// 				scriptPubKey: []byte{118,169,20,162,200,202,90,247,28,46,166,63,168,79,33,36,17,56,57,203,205,30,195,136,172},
// 			},
// 		},
// 	}

// 	genesisDifficulty, _ := DecodeDifficulty([]byte{0x1e, 0xff, 0xff, 0xff})
// 	//fmt.Printf("GEN DIFF %X \n", genesisDifficulty.targetBytes)
// 	genesisDifficulty = genesisDifficulty.Div(big.NewInt(diff))
// 	//fmt.Printf("GEN DIFF %X \n", genesisDifficulty.targetBytes)
// 	//genesisDifficulty, _ := DecodeDifficulty([]byte{0x1e, 0x19, 0xFF, 0xFF})
// 	// fmt.Printf("%X \n", genesisDifficulty.targetBytes)
// 	merkleRoot := MerkleRoot([]Transaction{coinbaseTx})
// 	bh := BlockHeader{
// 		PrevBlockHash: make([]byte, 32),
// 		merkleRoot: merkleRoot,
// 		timestamp: 1631303394,
// 		DifficultyTarget: genesisDifficulty,
// 		Nonce: 0,
// 	}

// 	//fmt.Printf("Merkle root %x \n", merkleRoot)

// 	return Block{
// 		Header: &bh,
// 		txs: []Transaction{coinbaseTx},
// 	}
// }
