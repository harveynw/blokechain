package chain

import (
	"fmt"
	"github.com/harveynw/blokechain/internal/data"
	"github.com/harveynw/blokechain/internal/script"
)

// Test tests this package
func Test() {
	// Sign transaction
	secretKey, _ := data.RandomKeyPair()

	// Recipient
	_, pk2 := data.RandomKeyPair()

	txIn := TransactionInput{
		prevTransaction: []byte("245e2d1f87415836cbb7b0bc84e40f4ca1d2a812be0eda381f02fb2224b4ad69"),
		prevIndex: 0,
		prevTransactionPubKey: []byte{0x01}, // FAKE LOCKING SCRIPT
		scriptSig: []byte{},
	}

	txOut := TransactionOutput{
		amount: 1000,
		scriptPubKey: script.P2PKH(pk2.HashEncode()).Encode(),
	}

	fmt.Printf("TX-IN %x \n", txIn.Encode())
	fmt.Printf("TX-OT %x \n", txOut.Encode())

	tx := Transaction{
		version: 1,
		txIn: []TransactionInput{txIn},
		txOut: []TransactionOutput{txOut},
	}

	txEncoded := tx.Encode(0)
	scriptSigComputed := data.SignMessage(secretKey, txEncoded)

	// Now we fill in scriptSig
	tx.txIn[0].scriptSig = scriptSigComputed.Encode()
	fmt.Println("Size of script sig", len(tx.txIn[0].scriptSig))

	txFullEncoded := tx.Encode(-1)
	fmt.Printf("TX: %x \n", txFullEncoded)

	txRecovered := DecodeTransaction(txFullEncoded)
	fmt.Println("TX recovered")
	fmt.Println(txRecovered)
}