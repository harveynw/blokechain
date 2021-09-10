package chain

import (
	"bytes"
	"fmt"
	"testing"
)

// TestDifficultyExamples checks that a selection of difficulties are encoded and decoded correctly
func TestDifficultyExamples(t *testing.T) {
	checkDecode([]byte{0x19, 0x03, 0xa3, 0x0c}, "0000000000000003A30C00000000000000000000000000000000000000000000", t)
	checkEncodeDecodeEquivalence([]byte{0x19, 0x03, 0xa3, 0x0c}, t)

	checkDecode([]byte{0x1b, 0x04, 0x04, 0xcb}, "00000000000404CB000000000000000000000000000000000000000000000000", t)
	checkEncodeDecodeEquivalence([]byte{0x1b, 0x04, 0x04, 0xcb}, t)

	checkDecode([]byte{0x1d, 0x00, 0xff, 0xff}, "00000000FFFF0000000000000000000000000000000000000000000000000000", t)
	checkEncodeDecodeEquivalence([]byte{0x1d, 0x00, 0xff, 0xff}, t)
}

func checkDecode(encoded []byte, targetHex string, t *testing.T) {
	diff, err := DecodeDifficulty(encoded)
	if err != nil {
		t.Fatalf("Failed decode %s", err)
	}

	outputHex := fmt.Sprintf("%X", diff.target.FillBytes(make([]byte, 32)))

	if outputHex != targetHex {
		t.Fatalf("Failed equivalence of decode\n %s == %s", targetHex, outputHex)
	}
}

func checkEncodeDecodeEquivalence(encoded []byte, t *testing.T) {
	diff, err := DecodeDifficulty(encoded)
	if err != nil {
		t.Fatalf("Failed decode %s", err)
	}
	if !bytes.Equal(encoded, diff.Encode()) {
		t.Fatalf("Failed equivalence of decode\n %s == %s", encoded, diff.Encode())
	}
}