package cryptography

import (
	"errors"
	"bytes"
	"math/big"
)

func fromPoint(p point) PublicKey {
	return PublicKey{p: p}
}

// HashEncode generates unique hash of public key
func (pk PublicKey) HashEncode() []byte {
	return Hash160(pk.EncodeCompressed())
}

// ToAddress gives a compressed public key address
func (pk PublicKey) ToAddress() string {
	pkHash := pk.HashEncode()
	verPkHash := append([]byte("\x00"), pkHash...) // Main Net

	checksumHash := Hash160(verPkHash)

	bytesAddress := append(verPkHash, checksumHash[:4]...)

	return Base58Encode(bytesAddress)
}

// Encode gives full (x, y) coordinates uncompressed for a public key
func (pk PublicKey) Encode() []byte {
	xBuf, yBuf := make([]byte, 32), make([]byte, 32)
	xBuf, yBuf = pk.p.x.FillBytes(xBuf), pk.p.y.FillBytes(yBuf)

	encoded := make([]byte, 0, 1 + 32 + 32)
	encoded = append(encoded, 0x04)
	encoded = append(append(encoded, xBuf...), yBuf...)

	return encoded
}

// EncodeCompressed gives x and sign of y for a public key
func (pk PublicKey) EncodeCompressed() []byte {
	var prefix []byte
	// Only need sign of y
	remainder := new(big.Int).Mod(&pk.p.y, big.NewInt(2))
	if remainder.Cmp(big.NewInt(0)) == 0 {
		prefix = []byte{0x02}
	} else {
		prefix = []byte{0x03}
	}

	return append(prefix, pk.p.x.Bytes()...)
}

// DecodePublicKey returns public key object from uncompressed format
func DecodePublicKey(b []byte) (PublicKey, error) {
	if len(b) != 1 + 32 + 32 || b[0] != 0x04 {
		return *new(PublicKey) , errors.New("Invalid format")
	}

	x := new(big.Int).SetBytes(b[1:33])
	y := new(big.Int).SetBytes(b[33:65])

	point := point{curve: &secp256k1, x: *x, y: *y}
	return PublicKey{p: point}, nil
}

// DecodePublicKeyCompressed returns public key object from the compressed format
func DecodePublicKeyCompressed(b []byte) (PublicKey, error) {
	if len(b) != 1 + 32 {
		return *new(PublicKey) , errors.New("Invalid format")
	}
	prefix := b[0]

	x := new(big.Int)
	var y *big.Int
	x.SetBytes(b[1:33])

	if prefix == 0x02 {
		// Even
		y, _ = YfromX(x)
	} else if prefix == 0x03 {
		// Odd
		_, y = YfromX(x)
	} else {
		return *new(PublicKey) , errors.New("Invalid format (sign)")
	}

	point := point{curve: &secp256k1, x: *x, y: *y}
	return PublicKey{p: point}, nil
}

// Encode signature using DER format
func (sig Signature) Encode() []byte {
	intEncode := func (i *big.Int) []byte {
		buf := make([]byte, 32)
		buf = i.FillBytes(buf)
		return append([]byte{0x02, byte(len(buf))}, buf...)
	}

	contents := append(intEncode(sig.r), intEncode(sig.s)...)
	return append([]byte{0x30, byte(len(contents))}, contents...)
}

// DecodeSignature recovers Signature from DER encoding
func DecodeSignature(b []byte) (Signature, error) {
	// SIGHASH_SINGLE only
	if b[0] != 0x30 {
		return *new(Signature) , errors.New("Invalid format")
	}

	contentLen := int(b[1])
	
	rStart, rHeader, rLen := 4, int(b[2]), int(b[3])
	rBytes := b[rStart:rStart+rLen]

	sStart, sHeader, sLen := rStart+rLen+2, int(b[rStart+rLen]), int(b[rStart+rLen+1])
	sBytes := b[sStart:sStart+sLen]


	if b[rHeader] != 0x02 || b[sHeader] != 0x02 || contentLen != 4 + rLen + sLen {
		return *new(Signature) , errors.New("Invalid r, s format")
	}

	r, s := new(big.Int).SetBytes(rBytes), new(big.Int).SetBytes(sBytes)

	return Signature{r: r, s: s}, nil
}

// Base58Encode Encodes arbitrary bytes
func Base58Encode(b []byte) string {
	x := new(big.Int).SetBytes(b)
	base58alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	base := big.NewInt(58)

	output := *new([]byte)

	for x.Cmp(big.NewInt(0)) != 0 {
		i := new(big.Int)
		x.DivMod(x, base, i)
		output = append(output, base58alphabet[i.Int64()])
	}

	// Leading zeros
	trimmed := bytes.TrimLeft(b, "\x00")
	pad := bytes.Repeat([]byte{base58alphabet[0]}, len(b) - len(trimmed))

	return string(append(pad, reverseBytes(output)...))
}

func reverseBytes(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

// EncodeInt returns int as a big-endian byte array
func EncodeInt(i int, nbytes int) []byte {
	buf := make([]byte, nbytes)
	buf = big.NewInt(int64(i)).FillBytes(buf)

	// if littleEndian {
	// 	buf = reverseBytes(buf)
	// }

	return buf
}

// EncodeVarInt returns big.Int as a variable length big-endian byte array
func EncodeVarInt(i int) []byte {
	if i < 0xfd {
		return []byte{byte(i)}
	} else if i < 0x10000 {
		return append([]byte{0xfd}, EncodeInt(i, 2)...)
	} else if i < 0x100000000 {
		return append([]byte{0xfe}, EncodeInt(i, 4)...)
	} else {
		return []byte{0xff}
	}
}

// DecodeInt returns int from big-endian encoded byte slice
func DecodeInt(b []byte) int {
	z := new(big.Int)
	z.SetBytes(b)
	return int(z.Int64())
}

// DecodeNextVarInt returns the next variable length integer encoded in b and the rest of b
func DecodeNextVarInt(b []byte) (int, []byte) {
	enc := b[0]
	if enc < 0xfd {
		return int(enc), b[1:]
	} else if enc == 0xfd {
		return DecodeInt(b[1:3]), b[3:]
	} else if enc == 0xfe {
		return DecodeInt(b[1:5]), b[5:]
	} else if enc == 0xff {
		return 0, b[1:]
	}
	return -1, b
}