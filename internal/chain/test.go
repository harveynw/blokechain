package chain

import (
	"fmt"
	"github.com/harveynw/blokechain/internal/data"
	"github.com/harveynw/blokechain/internal/script"
)

// Test tests this package
func Test() {
	prevHash := data.DoubleHash([]byte("Test"), false)

	// Sign transaction
	secretKey, pubKey := data.RandomKeyPair()
	// secretKey, pubKey := data.RandomKeyPair()

	// Recipient
	_, pk2 := data.RandomKeyPair()

	// Adversary
	// randomSecretKey, _ := data.RandomKeyPair()

	txIn := TransactionInput{
		prevTransaction: prevHash,
		prevIndex: 0,
		prevTransactionPubKey: script.P2PKH(pubKey.HashEncode()).Encode(), // LOCKING SCRIPT
		scriptSig: []byte{},
	}

	txOut := TransactionOutput{
		amount: 1000,
		scriptPubKey: script.P2PKH(pk2.HashEncode()).Encode(),
	}

	// fmt.Printf("TX-IN %x \n", txIn.Encode())
	// fmt.Printf("TX-OT %x \n", txOut.Encode())

	tx := Transaction{
		version: 1,
		txIn: []TransactionInput{txIn},
		txOut: []TransactionOutput{txOut},
	}

	txEncoded := tx.Encode(0)

	// Now we fill in scriptSig
	newScriptSig := script.NewScript()
	// newScriptSig.AppendData(data.SignMessage(randomSecretKey, txEncoded).Encode())
	newScriptSig.AppendData(data.SignMessage(secretKey, txEncoded).Encode())
	newScriptSig.AppendData(pubKey.EncodeCompressed())
	tx.txIn[0].scriptSig = newScriptSig.Encode()

	// fmt.Println("Size of script sig", len(newScriptSig.Encode()))

	// txFullEncoded := tx.Encode(-1)
	// fmt.Printf("TX: %x \n", txFullEncoded)

	// txRecovered := DecodeTransaction(txFullEncoded)
	// fmt.Println("TX recovered")
	// fmt.Println("Amount", txRecovered.txOut[0].amount)
	// fmt.Printf("Pubkey %x \n", txRecovered.txOut[0].scriptPubKey)
	// fmt.Printf("Before %x \n", script.P2PKH(pk2.HashEncode()).Encode())

	// fmt.Printf("Sig %x \n", txRecovered.txIn[0].scriptSig)
	// fmt.Printf("Sig before %x \n", newScriptSig.Encode())

	// fmt.Println("***NEW***")

	valid, err := tx.Verify()
	fmt.Println(valid, err)
}