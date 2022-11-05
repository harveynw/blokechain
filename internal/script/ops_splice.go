package script


func OP_CAT(vm *VM) bool {
	return false // Disabled
}

func OP_SUBSTR(vm *VM) bool {
	return false // Disabled
}

func OP_LEFT(vm *VM) bool {
	return false // Disabled
}

func OP_RIGHT(vm *VM) bool {
	return false // Disabled
}

func OP_SIZE(vm *VM) bool {
	if len(vm.Stack) < 1 {
		return false
	}

	l := int64(len(vm.Stack[0]))
	err, b := encodeInt(l)
	if err {
		return false
	}
	
	vm.Push(b, false)
	return true
}