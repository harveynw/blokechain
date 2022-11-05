package script

import (
	"bytes"
)

func OP_INVERT(vm *VM) bool {
	return false // Disabled
}

func OP_AND(vm *VM) bool {
	return false // Disabled
}

func OP_OR(vm *VM) bool {
	return false // Disabled
}

func OP_XOR(vm *VM) bool {
	return false // Disabled
}

func OP_EQUAL(vm *VM) bool {
	if len(vm.Stack) < 2 {
		return false
	}
	_, x1 := vm.Pop(false)
	_, x2 := vm.Pop(false)

	var v byte = 0x00
	if bytes.Compare(x1, x2) == 0 {
		v = 0x01
	}
	vm.Push([]byte{v}, false)
	return true
}

// OP_EQUALVERIFY Checks whether the top two elements of the stack are equal and then executes OP_VERIFY
func OP_EQUALVERIFY(vm *VM) bool {
	if len(vm.Stack) < 2 {
		return false
	}
	_, x1 := vm.Pop(false)
	_, x2 := vm.Pop(false)

	return bytes.Compare(x1, x2) == 0
}
