package chain

import (
	"fmt"
	"bytes"
	"github.com/harveynw/blokechain/internal/data"
)

// Transaction data structure containing multiple inputs and outputs
type Transaction struct {
	version int
	txIn []TransactionInput
	txOut []TransactionOutput
}

// TransactionInput data structure specifying an UXTO to be spent and unlocking script
type TransactionInput struct {
	prevTransaction []byte // Transaction hash containing UXTO
	prevIndex int // Select UXTO by index
	prevTransactionPubKey []byte // Previous script pubkey, needed for signature generation
	scriptSig []byte
}

// TransactionOutput data structure specifying spent coins and locking script
type TransactionOutput struct {
	amount uint64 // Satoshis
	scriptPubKey []byte // Unlocking script
}

// ID returns the transaction id SHA256(SHA256(transaction))
func(ts Transaction) ID() []byte {
	return data.DoubleHash(ts.Encode(-1), false)
}

// Encode transaction data structure using the protocol, signingIndex (if not -1) specifies the current input being using for signature generation
func (ts Transaction) Encode(signingIndex int) []byte {
	enc := make([]byte, 0)

	// Version
	enc = append(enc, 0x00, 0x00, 0x00, 0x01)

	// Input Counter
	enc = append(enc, data.EncodeVarInt(len(ts.txIn))...)
	// Inputs. If being encoded for ECDSA signature generation, replace publicKey
	for i, inTx := range ts.txIn {
		if signingIndex != -1 {
			if signingIndex == i {
				enc = append(enc, inTx.EncodeScriptSigOverride()...)
			} else {
				enc = append(enc, inTx.EncodeScriptSigEmpty()...)
			}
		} else {
			enc = append(enc, inTx.Encode()...)
		}
	}

	// Output Counter
	enc = append(enc, data.EncodeVarInt(len(ts.txOut))...)
	// Outputs
	for _, outTx := range ts.txOut {
		enc = append(enc, outTx.Encode()...)
	}

	// Locktime (ignore)
	enc = append(enc, 0x00, 0x00, 0x00, 0x00)

	return enc
}

// DecodeTransaction recovers Transaction according to the protocol
func DecodeTransaction(b []byte) Transaction {
	fmt.Println("BEGIN DECODE TX")
	compatibleVersion := []byte{0x00, 0x00, 0x00, 0x01}
	var version []byte
	version, b = b[0:4], b[4:]

	if bytes.Compare(version, compatibleVersion) != 0 {
		panic("Incompatible version")
	}

	inputCounter, b := data.DecodeNextVarInt(b)
	fmt.Println("INPUT COUNTER IS ", inputCounter)
	txIn := make([]TransactionInput, 0)
	for i := 0; i < inputCounter; i++ {
		var tx TransactionInput
		tx, b = DecodeNextTransactionInput(b)
		txIn = append(txIn, tx)
	}

	outputCounter, b := data.DecodeNextVarInt(b)
	fmt.Println("OUTPUT COUNTER IS ", outputCounter)
	txOut := make([]TransactionOutput, 0)
	for i := 0; i < outputCounter; i++ {
		var tx TransactionOutput
		tx, b = DecodeNextTransactionOutput(b)
		txOut = append(txOut, tx)
	}

	return Transaction{version: data.DecodeInt(version), txIn: txIn, txOut: txOut}
}

// Encode transaction input using protocol
func (in TransactionInput) Encode() []byte {
	enc := make([]byte, 0)

	// Previous transaction (32 bytes) + Output index (4 bytes)
	enc = append(enc, in.prevTransaction...)
	enc = append(enc, data.EncodeInt(in.prevIndex, 4)...)
	
	// Unlocking script size VarInt
	enc = append(enc, data.EncodeVarInt(len(in.scriptSig))...)

	// Unlocking script
	enc = append(enc, in.scriptSig...)

	// Sequence number (defunct)
	enc = append(enc, 0xFF, 0xFF, 0xFF, 0xFF)

	return enc
}

// EncodeScriptSigOverride encodes the transaction input, replacing the scriptSig with the previous tx pubKey (required for signature verification)
func (in TransactionInput) EncodeScriptSigOverride() []byte {
	in.scriptSig = in.prevTransactionPubKey
	return in.Encode()
}

// EncodeScriptSigEmpty encodes the transaction input, replacing the scriptSig an empty slice (required for signature verification)
func (in TransactionInput) EncodeScriptSigEmpty() []byte {
	in.scriptSig = make([]byte, 0)
	return in.Encode()
}

// DecodeNextTransactionInput recovers TransactionInput according to the protocol and returns rest of data
func DecodeNextTransactionInput(b []byte) (TransactionInput, []byte) {
	prevTransaction := b[0:32]
	prevIndex := data.DecodeInt(b[32:36])

	fmt.Printf("scriptSigSizeData %x \n", b[36:])
	scriptSigSize, b := data.DecodeNextVarInt(b[36:])
	fmt.Println("INPUT scriptSig size is ", scriptSigSize)
	scriptSig := b[0:scriptSigSize]

	return TransactionInput{prevTransaction: prevTransaction, prevIndex: prevIndex, scriptSig: scriptSig}, b[scriptSigSize+4:]
}

// Encode transaction output using protocol
func (out TransactionOutput) Encode() []byte {
	enc := make([]byte, 0)

	// Amount in satoshis
	enc = append(enc, data.EncodeInt(int(out.amount), 8)...)

	// Locking script size
	enc = append(enc, data.EncodeVarInt(len(out.scriptPubKey))...)

	// Locking script
	enc = append(enc, out.scriptPubKey...)

	return enc
}

// DecodeNextTransactionOutput recovers TransactionOutput according to the protocol and returns rest of data
func DecodeNextTransactionOutput(b []byte) (TransactionOutput, []byte) {
	amount := uint64(data.DecodeInt(b[0:8]))
	scriptPubKeySize, b := data.DecodeNextVarInt(b[8:])

	return TransactionOutput{amount: amount, scriptPubKey: b[0:scriptPubKeySize]}, b[scriptPubKeySize:]
}