package script

import (
	"encoding/binary"
)

func isTruthy(b []byte) bool {
	// Anything not zero
	return !isZero(b)
}

func isFalse(b []byte) bool {
	// Equivalent to isZero
	return isZero(b)
}

func isZero(b []byte) bool {
	// Positive zero
	if len(b) == 0 {
		return true
	}

	// Test for arbitrary length negative/positive zero
	if !(b[0] == 0x80 || b[0] == 0x00) {
		return false
	}
	for _, v := range b[1:] {
		if v != 0x00 {
			return false
		}
	}
	return true
}

func decodeInt(b []byte) (err bool, val int64) {
	l := len(b)

	// Arithmetic on more than 4 bytes input forbidden
	if l > 4 {
		return true, 0
	}
	if l == 0 {
		return false, 0
	}

	// Determine sign and remove sign bit
	t := b[l-1]
	isNeg := int64((t >> 7) & 0x01)
	b[l-1] = (t << 1) >> 1
	//fmt.Printf("isNeg = %v, urep = %v \n", isNeg, b)

	// Right pad to 8 bytes total
	for i := 8; i > l; i-- {
		b = append(b, 0x00)
	}

	//fmt.Printf("Padded %v \n", b)

	// Cast to (signed) int64
	val = int64(binary.LittleEndian.Uint64(b))
	val = (1 - 2*isNeg) * val

	return false, val
}

func encodeInt(i int64) (err bool, b []byte) {
	if i == 0 {
		return false, []byte{0x00}
	}

	var isNeg bool = i < 0
	if isNeg {
		i *= -1
	}

	// To little-endian byte slice
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))

	// Trim leading zero bytes
	for i := 7; i >= 0; i-- {
		if b[i] == 0x00 {
			b = b[0:i]
		} else {
			break
		}
	}

	if isNeg && len(b) > 0 {
		if b[len(b)-1] >= 0x80 {
			// Append new negative bit byte
			b = append(b, 0x80)
		} else {
			// Flip negative bit
			b[len(b)-1] = b[len(b)-1] | 0x80
		}
	}

	// Script spec
	if len(b) > 5 {
		return true, nil
	}
	return false, b
}