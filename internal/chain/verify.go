package chain

import (
	"errors"
	"github.com/harveynw/blokechain/internal/script"
)

// ErrPubKeyMissing When the previous unlocking script (public key) to verify the transaction is missing
var ErrPubKeyMissing = errors.New("Previous transaction public key is missing")
// ErrScriptSigInvalid When the provided signature does not match the original public key
var ErrScriptSigInvalid = errors.New("Unlock script failed")

// Verify by executing the unlock + locking script and checking input >= output
func (ts Transaction) Verify() (bool, error) {
	for i, txIn := range ts.txIn {
		sigEncoding := ts.Encode(i) 
		valid, err := verifyTransactionInput(sigEncoding, txIn)

		if !(valid && err == nil) {
			return false, err
		}
	}

	// TODO input >= output
	return true, nil
}

func verifyTransactionInput(sigEncoding []byte, txIn TransactionInput) (bool, error) {
	if len(txIn.prevTransactionPubKey) == 0 {
		return false, ErrPubKeyMissing
	}

	unlock, lock := script.DecodeScript(txIn.scriptSig), script.DecodeScript(txIn.prevTransactionPubKey)
	unlock.Concat(lock)
	unlock.Print()
	result := unlock.Execute(sigEncoding)

	if result {
		return true, nil
	}
	return false, ErrScriptSigInvalid
}