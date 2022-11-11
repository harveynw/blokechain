package chain

type BlockHeader struct {
	hashPrevBlock []byte
	hashMerkleRoot []byte
	Time uint32
	DifficultyTarget Difficulty
	Nonce uint32
}

func (bh BlockHeader) Encode() []byte {
	enc := make([]byte, 0)

	// Version (arbitrary)
	enc = append(enc, 0x02, 0x00, 0x00, 0x01)
	
	// Previous Block Hash + Merkle Root (32 Bytes)
	enc = append(enc, bh.hashPrevBlock[0:32]...)
	enc = append(enc, bh.hashMerkleRoot[0:32]...)

	// Timestamp (seconds from Unix Epoch, 4 Bytes)
	enc = append(enc, EncodeInt(int64(bh.Time), 4)...)

	// Difficulty Target + Nonce
	enc = append(enc, bh.DifficultyTarget.Encode()...)
	enc = append(enc, EncodeInt(int64(bh.Nonce), 4)...)

	return enc
}

func DecodeNextBlockHeader(b []byte) (BlockHeader, []byte) {
	b = b[4:] // Discard Version

	// Block headers are fixed length
	hashPrevBlock, b := b[0:32], b[32:]
	hashMerkleRoot, b := b[0:32], b[32:]
	timestampBytes, b := b[0:4], b[4:]
	difficultyBytes, b := b[0:4], b[4:]
	nonceBytes, b := b[0:4], b[4:]

	difficulty, err := DecodeDifficulty(difficultyBytes)
	if err != nil {
		panic("Malformed difficulty in block header")
	}

	return BlockHeader{
		hashPrevBlock: hashPrevBlock,
		hashMerkleRoot: hashMerkleRoot,
		Time: uint32(DecodeInt(timestampBytes)),
		DifficultyTarget: difficulty,
		Nonce: uint32(DecodeInt(nonceBytes)),
	}, b
}

