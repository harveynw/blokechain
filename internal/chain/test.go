package chain

// import (
// 	"fmt"
// 	"github.com/harveynw/blokechain/internal/data"
// 	"github.com/harveynw/blokechain/internal/script"
// )

// // Test tests this package
// func Test() {
// 	prevHash := data.DoubleHash([]byte("Test"), false)

// 	// Sign transaction
// 	secretKey, pubKey := data.RandomKeyPair()

// 	// Recipient
// 	_, pk2 := data.RandomKeyPair()


// 	txIn := TransactionInput{
// 		prevTransaction: prevHash,
// 		prevIndex: 0,
// 		prevTransactionPubKey: script.P2PKH(pubKey.HashEncode()).Encode(), // LOCKING SCRIPT
// 		scriptSig: []byte{},
// 	}

// 	txOut := TransactionOutput{
// 		amount: 1000,
// 		scriptPubKey: script.P2PKH(pk2.HashEncode()).Encode(),
// 	}

// 	tx := Transaction{
// 		version: 1,
// 		txIn: []TransactionInput{txIn},
// 		txOut: []TransactionOutput{txOut},
// 	}

// 	txEncoded := tx.Encode(0)

// 	// Now we fill in scriptSig
// 	newScriptSig := script.NewScript()
// 	newScriptSig.AppendData(data.SignMessage(secretKey, txEncoded).Encode())
// 	newScriptSig.AppendData(pubKey.EncodeCompressed())
// 	tx.txIn[0].scriptSig = newScriptSig.Encode()

// 	valid, err := tx.Verify()
// 	fmt.Println(valid, err)
// }