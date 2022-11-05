package script

import (
	"github.com/harveynw/blokechain/internal/cryptography"
)


// OP_RIPEMD160 Hashes the input using RIPEMD160 
func OP_RIPEMD160(vm *VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.RIPEMD160(value), false)
	return true
}

// OP_SHA1 Hashes the input using SHA1
func OP_SHA1(vm *VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.SHA1(value), false)
	return true
}

// OP_SHA256 Hashes the input using SHA256
func OP_SHA256(vm *VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.SHA256(value), false)
	return true
}

// OP_HASH160 Hashes the input using Hash160
func OP_HASH160(vm *VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.Hash160(value), false)
	return true
}

// OP_HASH256 Hashes the input using Hash256
func OP_HASH256(vm *VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.Hash256(value), false)
	return true
}

// OP_CODESEPERATOR Does nothing
func OP_CODESEPERATOR(vm *VM) bool {
	// TODO
	return true
}

// OP_CHECKSIG Pushes true/false, depending on whether the pub key and signature are valid for the transaction
func OP_CHECKSIG(vm *VM) bool {
	err1, pubKeyBytes := vm.Pop(false)
	err2, sigBytes := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	sig, err3 := cryptography.DecodeSignature(sigBytes)
	if err3 != nil {
		return false
	}
	pubKey, err4 := cryptography.DecodePublicKeyCompressed(pubKeyBytes)
	if err4 != nil {
		return false
	}

	if sig.VerifySignature(pubKey, vm.Transaction) {
		vm.Push([]byte{0x01}, false) // Truthy
	} else {
		vm.Push([]byte{}, false) // False
	}

	return true
}

func OP_CHECKSIGVERIFY(vm *VM) bool {
	err1 := OP_CHECKSIG(vm)
	if err1 {
		return false
	}

	err2 := OP_VERIFY(vm)
	return !err2
}

func OP_CHECKMULTISIG(vm *VM) bool {
	// Public Keys
	err_m, m_b := vm.Pop(false)
	if err_m {
		return false
	}
	err_dec_m, m := decodeInt(m_b)
	if err_dec_m {
		return false
	}
	pks := make([]cryptography.PublicKey, 0)
	for i := int64(0); i < m; i++ {
		err_pk, pk_b := vm.Pop(false)
		if err_pk {
			return false
		}
		pubKey, pk_dec_err := cryptography.DecodePublicKeyCompressed(pk_b)
		if pk_dec_err != nil {
			return false
		}
		pks = append(pks, pubKey)
	}

	// Signatures
	err_n, n_b := vm.Pop(false)
	if err_n {
		return false
	}
	err_dec_n, n := decodeInt(n_b)
	if err_dec_n {
		return false
	}
	sigs := make([]cryptography.Signature, 0)
	for i := int64(0); i < n; i++ {
		err_sig, sig_b := vm.Pop(false)
		if err_sig {
			return false
		}
		signature, sig_dec_err := cryptography.DecodeSignature(sig_b)
		if sig_dec_err != nil {
			return false
		}
		sigs = append(sigs, signature)
	}

	// BIP 147 Requirement
	err_bug, dummy := vm.Pop(false)
	if err_bug || !isZero(dummy) {
		return false
	}

	// Check signatures
	for _, signature := range sigs {
		pks_left := len(pks)
		for idx := 0; idx < pks_left; idx++ {
			if signature.VerifySignature(pks[idx], vm.Transaction) {
				pks = append(pks[:idx], pks[idx+1:]...)
				continue
			}
		}
		return false
	}

	return true
}

func OP_CHECKMULTISIGVERIFY(vm *VM) bool {
	err1 := OP_CHECKMULTISIG(vm)
	if err1 {
		return false
	}

	err2 := OP_VERIFY(vm)
	return !err2
}