package ops

import (
	"github.com/harveynw/blokechain/internal/script"
)

// OP_NOP
func OP_NOP(vm *script.VM) bool {
	return true
}

// OP_VERIFY Is top value of stack truthy
func OP_VERIFY(vm *script.VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	return script.IsTruthy(value)
}

// OP_RETURN Fails immediately
func OP_RETURN(vm *script.VM) bool {
	return false
}
