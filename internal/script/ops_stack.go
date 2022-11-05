package script


// OP_TOALTSTACK 
func OP_TOALTSTACK(vm *VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}

	vm.Push(value, true)
	return true
}

func OP_FROMALTSTACK(vm *VM) bool {
	err, value := vm.Pop(true)
	if err {
		return false
	}

	vm.Push(value, false)
	return true
}

func OP_IFDUP(vm *VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}

	if !isZero(value) {
		vm.Push(value, false)
		vm.Push(value, false)
	}

	return true
}

func OP_DEPTH(vm *VM) bool {
	depth := len(vm.Stack)
	
	err, b := encodeInt(int64(depth))
	if err {
		return false
	}

	vm.Push(b, false)
	return true
}

func OP_DROP(vm *VM) bool {
	err, _ := vm.Pop(false)
	return !err
}

// OP_DUP Duplicates top element of the stack
func OP_DUP(vm *VM) bool {
	err, value := vm.Pop(false)
	if err {
		return false
	}
	vm.Push(value, false)
	vm.Push(value, false)
	return true
}

func OP_NIP(vm *VM) bool {
	err1, val := vm.Pop(false)
	err2, _ := vm.Pop(false)

	if err1 || err2 {
		return false
	}

	vm.Push(val, false)
	return true
}

func OP_OVER(vm *VM) bool {
	err1, val1 := vm.Pop(false)
	err2, val2 := vm.Pop(false)

	if err1 || err2 {
		return false
	}

	vm.Push(val2, false)
	vm.Push(val1, false)
	vm.Push(val2, false)
	return true
}

func OP_PICK(vm *VM) bool {
	err1, nb := vm.Pop(false)
	if err1 {
		return false
	}
	err2, n := decodeInt(nb)
	if err2 {
		return false
	}

	idx := int(n)
	if idx >= len(vm.Stack) {
		return false
	}
	vm.Push(vm.Stack[idx], false)
	return true
}

func OP_ROLL(vm *VM) bool {
	err1, nb := vm.Pop(false)
	if err1 {
		return false
	}
	err2, n := decodeInt(nb)
	if err2 {
		return false
	}

	idx := int(n)
	if idx >= len(vm.Stack) {
		return false
	}

	v := vm.Stack[idx]
	vm.Stack = append(vm.Stack[:idx], vm.Stack[idx+1:]...)
	vm.Push(v, false)
	return true
}

func OP_ROT(vm *VM) bool {
	if len(vm.Stack) < 3 {
		return false
	}

	vm.Stack[0], vm.Stack[1], vm.Stack[2] = vm.Stack[2], vm.Stack[0], vm.Stack[1]
	return true
}

func OP_SWAP(vm *VM) bool {
	if len(vm.Stack) < 2 {
		return false
	}

	vm.Stack[0], vm.Stack[1] = vm.Stack[1], vm.Stack[0]
	return true
}

func OP_TUCK(vm *VM) bool {
	if len(vm.Stack) < 2 {
		return false
	}

	vm.Stack[0], vm.Stack[1] = vm.Stack[1], vm.Stack[0]
	vm.Push(vm.Stack[1], false)
	return true
}

func OP_2DROP(vm *VM) bool {
	if len(vm.Stack) < 2 {
		return false
	}

	_, _ = vm.Pop(false)
	_, _ = vm.Pop(false)
	return true
}

func OP_2DUP(vm *VM) bool {
	if len(vm.Stack) < 2 {
		return false
	}

	_, v2 := vm.Pop(false)
	_, v1 := vm.Pop(false)

	vm.Push(v1, false)
	vm.Push(v2, false)
	vm.Push(v1, false)
	vm.Push(v2, false)
	return true
}

func OP_3DUP(vm *VM) bool {
	if len(vm.Stack) < 3 {
		return false
	}

	_, v3 := vm.Pop(false)
	_, v2 := vm.Pop(false)
	_, v1 := vm.Pop(false)

	vm.Push(v1, false)
	vm.Push(v2, false)
	vm.Push(v3, false)
	vm.Push(v1, false)
	vm.Push(v2, false)
	vm.Push(v3, false)
	return true
}

func OP_2OVER(vm *VM) bool {
	if len(vm.Stack) < 4 {
		return false
	}

	v1, v2 := vm.Stack[2], vm.Stack[3]
	vm.Push(v2, false)
	vm.Push(v1, false)
	return true
}

func OP_2ROT(vm *VM) bool {
	if len(vm.Stack) < 6 {
		return false
	}

	x6, x5, x4, x3, x2, x1 := vm.Stack[0], vm.Stack[1], vm.Stack[2], vm.Stack[3], vm.Stack[4], vm.Stack[5]
	vm.Stack[0], vm.Stack[1], vm.Stack[2], vm.Stack[3], vm.Stack[4], vm.Stack[5] = x2, x1, x6, x5, x4, x3
	return true
}

func OP_2SWAP(vm *VM) bool {
	if len(vm.Stack) < 4 {
		return false
	}

	x4, x3, x2, x1 := vm.Stack[0], vm.Stack[1], vm.Stack[2], vm.Stack[3]
	vm.Stack[0], vm.Stack[1], vm.Stack[2], vm.Stack[3] = x2, x1, x4, x3
	return true
}