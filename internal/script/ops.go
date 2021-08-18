package script

import (
	"bytes"
	"github.com/harveynw/blokechain/internal/data"
)

// OP_DUP Duplicates top element of the stack
func OP_DUP(s *ExecutionStack) bool {
	err, value := s.Pop()
	if err {
		return false
	}
	s.Push(value)
	s.Push(value)
	return true
}

// OP_HASH160 Removes and adds a 20 byte hash of the top element of the stack 
func OP_HASH160(s *ExecutionStack) bool {
	err, value := s.Pop()
	if err {
		return false
	}
	s.Push(data.DoubleHash(value, true))
	return true
}

// OP_EQUALVERIFY Checks whether the top two elements of the stack are equal and then executes OP_VERIFY
func OP_EQUALVERIFY(s *ExecutionStack) bool {
	err1, val1 := s.Pop()
	err2, val2 := s.Pop()
	if err1 || err2 {
		return false
	}
	return bytes.Compare(val1, val2) == 0
}

// OP_CHECKSIG Pushes true/false, depending on whether the pub key and signature are valid for the transaction
func OP_CHECKSIG(s *ExecutionStack) bool {
	err1, pubKeyBytes := s.Pop()
	err2, sigBytes := s.Pop()
	if err1 || err2 {
		return false
	}

	sig, err3 := data.DecodeSignature(sigBytes)
	if err3 != nil {
		return false
	}
	pubKey, err4 := data.DecodePublicKeyCompressed(pubKeyBytes)
	if err4 != nil {
		return false
	}

	if sig.VerifySignature(pubKey, s.TransactionEncoded) {
		s.Push([]byte{0x01}) // Truthy
	} else {
		s.Push([]byte{}) // False
	}

	return true
}
