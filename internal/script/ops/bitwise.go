package ops

import (
	"bytes"
	"github.com/harveynw/blokechain/internal/script"
)

// OP_EQUALVERIFY Checks whether the top two elements of the stack are equal and then executes OP_VERIFY
func OP_EQUALVERIFY(vm *script.VM) bool {
	err1, val1 := vm.Pop(false)
	err2, val2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}
	return bytes.Compare(val1, val2) == 0
}
