package chain

import (
	"fmt"
	"testing"
	"bytes"
	"encoding/hex"
)

func TestDecodeGenesis(t *testing.T) {
	// Random tx
	// ee475443f1fbfff84ffba43ba092a70d291df233bd1428f3d09f7bd1a6054a1f
	txraw := "010000000110ee96aa946338cfd0b2ed0603259cfe2f5458c32ee4bd7b88b583769c6b046e010000006b483045022100e5e4749d539a163039769f52e1ebc8e6f62e39387d61e1a305bd722116cded6c022014924b745dd02194fe6b5cb8ac88ee8e9a2aede89e680dcea6169ea696e24d52012102b4b754609b46b5d09644c2161f1767b72b93847ce8154d795f95d31031a08aa2ffffffff028098f34c010000001976a914a134408afa258a50ed7a1d9817f26b63cc9002cc88ac8028bb13010000001976a914fec5b1145596b35f59f8be1daf169f375942143388ac00000000"
	transactionBytes, _ := hex.DecodeString(txraw)

	tx, _ := DecodeNextTransaction(transactionBytes)

	fmt.Println(tx)
}

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
