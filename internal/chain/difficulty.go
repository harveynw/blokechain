package chain

import (
	"bytes"
	"errors"
	"math/big"
	//"github.com/harveynw/blokechain/internal/params"
)

// Difficulty houses logic for computing, encoding and testing difficulty
type Difficulty struct {
	target *big.Int
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func leftPad(b []byte, l int) ([]byte, int) {
	if len(b) >= l {
		return b, 0
	}

	diff := l - len(b)
	return append(make([]byte, diff), b...), diff
}

// DecodeDifficulty recovers from mantissa-exponent format
func DecodeDifficulty(b []byte) (Difficulty, error) {
	if len(b) != 4 {
		return Difficulty{}, errors.New("Difficulty field wrong size")
	}
	coefficient := new(big.Int).SetBytes(b[1:4])
	exponent := uint(8*(int(b[0])-3))

	target := new(big.Int)
	target.Lsh(coefficient, exponent)

	return Difficulty{target: target}, nil
}

// Encode returns a 4 byte mantissa-exponent encoding of the difficulty
func (diff Difficulty) Encode() []byte {
	b := diff.target.FillBytes(make([]byte, 32))
	var mantissa int
	for i, digit := range b {
		if digit != 0x00 {
			mantissa = len(b) - i
			b = b[i:min(i+3, len(b))]
			break
		}
	}
	b = bytes.TrimRight(b, string(byte(0x00)))
	b, amountPadded := leftPad(b, 3) // Ensure its 3 bytes
	mantissa += amountPadded

	return append([]byte{byte(mantissa)}, b...)
}