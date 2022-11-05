package script

// OP_NOP
func OP_NOP(vm *VM) bool {
	return true
}

// OP_VERIFY Is top value of stack truthy
func OP_VERIFY(vm *VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	return isTruthy(value)
}

// OP_RETURN Fails immediately
func OP_RETURN(vm *VM) bool {
	return false
}
