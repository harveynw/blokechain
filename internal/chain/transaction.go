package chain

import (
	"bytes"
	"github.com/harveynw/blokechain/internal/cryptography"
)

// Transaction data structure containing multiple inputs and outputs
type Transaction struct {
	isSegwit bool
	txIn []TransactionInput
	txOut []TransactionOutput
	Witnesses []byte // todo
	lock_time Locktime
}

// TransactionInput data structure specifying an UXTO to be spent and unlocking script
type TransactionInput struct {
	prevTransaction []byte // Transaction hash containing UXTO
	prevIndex int64 // Select UXTO by index
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
	return cryptography.Hash256(ts.Encode(-1))
}

// Encode transaction data structure using the protocol, signingIndex (if not -1) specifies the current input being using for signature generation
func (ts Transaction) Encode(signingIndex int) []byte {
	enc := make([]byte, 0)

	// Version, always 1
	enc = append(enc, 0x00, 0x00, 0x00, 0x01)

	// If witness data present, else omitted
	if ts.isSegwit {
		enc = append(enc, 0x00, 0x01)
	}

	// Input Counter
	n_inputs := NewVarInt(len(ts.txIn))
	enc = append(enc, n_inputs.EncodeVarInt()...)
	// Inputs. If being encoded for ECDSA signature generation (signingIndex != -1) -> replace publicKey
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
	n_outputs := NewVarInt(len(ts.txIn))
	enc = append(enc, n_outputs.EncodeVarInt()...)
	// Outputs
	for _, outTx := range ts.txOut {
		enc = append(enc, outTx.Encode()...)
	}

	// Witness data
	if ts.isSegwit {
		// TODO
		panic("SegWit not implemented")
	}

	// Locktime
	enc = append(enc, ts.lock_time.Encode()...)

	return enc
}

// DecodeTransaction recovers Transaction according to the protocol
func DecodeNextTransaction(b []byte) (Transaction, []byte) {
	var version []byte
	version, b = b[0:4], b[4:]
	if bytes.Compare(version, []byte{0x00, 0x00, 0x00, 0x01}) != 0 {
		panic("Incompatible version")
	}

	isSegwit := false
	if bytes.Compare(b[0:2], []byte{0x00, 0x01}) == 0 {
		isSegwit, b = true, b[2:]
	}

	inputCounter, b := DecodeNextVarInt(b)
	txIn := make([]TransactionInput, 0)
	for i := 0; i < int(inputCounter.val); i++ {
		var tx TransactionInput
		tx, b = DecodeNextTransactionInput(b)
		txIn = append(txIn, tx)
	}

	outputCounter, b := DecodeNextVarInt(b)
	txOut := make([]TransactionOutput, 0)
	for i := 0; i < int(outputCounter.val); i++ {
		var tx TransactionOutput
		tx, b = DecodeNextTransactionOutput(b)
		txOut = append(txOut, tx)
	}

	if isSegwit {
		// TODO Witnesses
		panic("SegWit not implemented")
	}

	lock_time := DecodeLocktime(b[0:4])

	return Transaction{
		isSegwit: isSegwit,
		txIn: txIn,
		txOut: txOut,
		Witnesses: nil,
		lock_time: lock_time,
	}, b[4:]
}

// Encode transaction input using protocol
func (in TransactionInput) Encode() []byte {
	enc := make([]byte, 0)

	if len(in.prevTransaction) != 32 {
		panic("Invalid previous transaction hash, must be 32 bytes")
	}

	// Previous transaction (32 bytes) + Output index (4 bytes)
	enc = append(enc, in.prevTransaction...)
	enc = append(enc, EncodeInt(int64(in.prevIndex), 4)...)
	
	// Unlocking script size VarInt
	script_size := NewVarInt(len(in.scriptSig))
	enc = append(enc, script_size.EncodeVarInt()...)

	// Unlocking script
	enc = append(enc, in.scriptSig...)

	// Sequence number (not used)
	enc = append(enc, 0xFF, 0xFF, 0xFF, 0xFF)

	return enc
}

// EncodeScriptSigOverride encodes the transaction input, replacing the scriptSig with the previous tx pubKey (required for signature verification)
func (in TransactionInput) EncodeScriptSigOverride() []byte {
	in.scriptSig = in.prevTransactionPubKey
	return in.Encode()
}

// EncodeScriptSigEmpty encodes the transaction input, replacing the scriptSig with an empty slice (required for signature verification)
func (in TransactionInput) EncodeScriptSigEmpty() []byte {
	in.scriptSig = make([]byte, 0)
	return in.Encode()
}

// DecodeNextTransactionInput recovers TransactionInput according to the protocol and returns rest of data
func DecodeNextTransactionInput(b []byte) (TransactionInput, []byte) {
	prevTransaction := b[0:32]
	prevIndex := DecodeInt(b[32:36])

	scriptSigSizeVarInt, b := DecodeNextVarInt(b[36:])
	scriptSigSize := int(scriptSigSizeVarInt.val)
	scriptSig := b[0:scriptSigSize]

	if len(b[scriptSigSize:]) != 4 || bytes.Compare(b[scriptSigSize:], []byte{0xFF, 0xFF, 0xFF, 0xFF}) != 0 {
		panic("Expected 4 bytes for sequence_no")
	}

	return TransactionInput{prevTransaction: prevTransaction, prevIndex: prevIndex, scriptSig: scriptSig}, b[scriptSigSize+4:]
}

// Encode transaction output using protocol
func (out TransactionOutput) Encode() []byte {
	enc := make([]byte, 0)

	// Amount in satoshis
	enc = append(enc, EncodeInt(int64(out.amount), 8)...)

	// Locking script size
	script_size := NewVarInt(len(out.scriptPubKey))
	enc = append(enc, script_size.EncodeVarInt()...)

	// Locking script
	enc = append(enc, out.scriptPubKey...)

	return enc
}

// DecodeNextTransactionOutput recovers TransactionOutput according to the protocol and returns rest of data
func DecodeNextTransactionOutput(b []byte) (TransactionOutput, []byte) {
	amount := uint64(DecodeInt(b[0:8]))
	scriptPubKeySizeVarInt, b := DecodeNextVarInt(b[8:])
	scriptPubKeySize := int(scriptPubKeySizeVarInt.val)

	return TransactionOutput{amount: amount, scriptPubKey: b[0:scriptPubKeySize]}, b[scriptPubKeySize:]
}