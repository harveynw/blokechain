package ops

import (
	"github.com/harveynw/blokechain/internal/cryptography"
	"github.com/harveynw/blokechain/internal/script"
)


// OP_RIPEMD160 Hashes the input using RIPEMD160 
func OP_RIPEMD160(vm *script.VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.RIPEMD160(value), false)
	return true
}

// OP_SHA1 Hashes the input using SHA1
func OP_SHA1(vm *script.VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.SHA1(value), false)
	return true
}

// OP_SHA256 Hashes the input using SHA256
func OP_SHA256(vm *script.VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.SHA256(value), false)
	return true
}

// OP_HASH160 Hashes the input using Hash160
func OP_HASH160(vm *script.VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.Hash160(value), false)
	return true
}

// OP_HASH256 Hashes the input using Hash256
func OP_HASH256(vm *script.VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(cryptography.Hash256(value), false)
	return true
}

// OP_CODESEPERATOR Does nothing
func OP_CODESEPERATOR(vm *script.VM) bool {
	// TODO
	return true
}

// OP_CHECKSIG Pushes true/false, depending on whether the pub key and signature are valid for the transaction
func OP_CHECKSIG(vm *script.VM) bool {
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

func OP_CHECKSIGVERIFY(vm *script.VM) bool {
	err1 := OP_CHECKSIG(vm)
	if err1 {
		return false
	}

	err2 := OP_VERIFY(vm)
	return !err2
}

func OP_CHECKMULTISIG(vm *script.VM) bool {
	// TODO
	return true
}

func OP_CHECKMULTISIGVERIFY(vm *script.VM) bool {
	err1 := OP_CHECKMULTISIG(vm)
	if err1 {
		return false
	}

	err2 := OP_VERIFY(vm)
	return !err2
}