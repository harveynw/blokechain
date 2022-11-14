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