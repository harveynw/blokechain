package chain

import (
	"math/big"
)

type VarInt struct {
	val int64
}

func NewVarInt(val int) *VarInt {
	return &VarInt{val: int64(val)}
}

// EncodeBytes encodes as [big-endian...]
func (vi *VarInt) EncodeBytes(nbytes int) []byte {
	buf := make([]byte, nbytes)
	buf = big.NewInt(vi.val).FillBytes(buf)
	return buf
}

// EncodeVarInt encodes as [nbytes, big-endian...]
func (vi *VarInt) EncodeVarInt() []byte {
	if vi.val < 0xfd {
		return []byte{byte(vi.val)}
	} else if vi.val < 0x10000 {
		return append([]byte{0xfd}, vi.EncodeBytes(2)...)
	} else if vi.val < 0x100000000 {
		return append([]byte{0xfe}, vi.EncodeBytes(4)...)
	} else {
		return []byte{0xff}
	}
}

// Decodes [nbytes, big-endian...] into VarInt struct, returning rest of b
func DecodeNextVarInt(b []byte) (*VarInt, []byte) {
	nBytes := b[0]
	if nBytes < 0xfd {
		return &VarInt{val: int64(nBytes)}, b[1:]
	} else if nBytes == 0xfd {
		return &VarInt{val: DecodeInt(b[1:3])}, b[3:]
	} else if nBytes == 0xfe {
		return &VarInt{val: DecodeInt(b[1:5])}, b[5:]
	} else if nBytes == 0xff {
		// Unsupported past 64-bit
		return nil, b
	}
	return nil, b
}

// EncodeInt encodes i as [big-endian...]
func EncodeInt(i int64, nbytes int) []byte {
	buf := make([]byte, nbytes)
	buf = big.NewInt(i).FillBytes(buf)
	return buf
}

// DecodeInt returns int from big-endian encoded byte slice
func DecodeInt(b []byte) int64 {
	z := new(big.Int)
	z.SetBytes(b)
	return z.Int64()
}

func reverseBytes(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}  