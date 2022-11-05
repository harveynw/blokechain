package script

func OP_0(vm *VM) bool {
	vm.Push([]byte{}, false)
	return true
}

func OP_1NEGATE(vm *VM) bool {
	vm.Push([]byte{0x81}, false)
	return true
}

func OP_TRUE(vm *VM) bool {
	vm.Push([]byte{0x01}, false)
	return true
}

func OP_N(n int) func(*VM) bool {
	return func(vm *VM) bool {
		vm.Push([]byte{byte(n)}, false)
		return true
	}
}