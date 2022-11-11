package chain

import (
	"bytes"
	"errors"
	"math/big"
)

// Difficulty houses logic for computing, encoding and testing difficulty
type Difficulty struct {
	target *big.Int
	targetBytes []byte
}

func leftPad(b []byte, l int) ([]byte, int) {
	if len(b) >= l {
		return b, 0
	}

	diff := l - len(b)
	return append(make([]byte, diff), b...), diff
}

func Compare(a, b []byte) int {
	if len(a) != len(b) {
		return -2
	}

	for i, val := range a {
		if val < b[i] {
			return -1
		}
		if val > b[i] {
			return +1
		}
	}
	return 0
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

	targetBytes := target.FillBytes(make([]byte, 32))

	return Difficulty{target: target, targetBytes: targetBytes}, nil
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

// IsSolution tests whether a block hash is below the difficulty target
func (diff Difficulty) IsSolution(hash []byte) bool {
	if Compare(hash, diff.targetBytes) == -1 {
		//fmt.Println("FOUND!")
		//fmt.Printf("%X < %X \n", hash, diff.targetBytes)
		return true
	}
	return false
	//return compare(hash, diff.targetBytes) == -1
}

// Mul multiples the difficulty by constant
func (diff Difficulty) Mul(c *big.Int) Difficulty {
	diff.target.Mul(diff.target, c)
	diff.target.FillBytes(diff.targetBytes)
	return diff
}

// Div divides the difficulty by constant
func (diff Difficulty) Div(c *big.Int) Difficulty {
	// fmt.Println("Diff-b", diff.target)
	diff.target.Div(diff.target, c)
	diff.target.FillBytes(diff.targetBytes)
	// fmt.Println("Diff-a", diff.target)
	return diff
}