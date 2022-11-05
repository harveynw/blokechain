package script


func OP_1ADD(vm *VM) bool {
	err1, v := vm.Pop(false)
	if err1 {
		return false
	}

	err2, i := decodeInt(v)
	if err2 {
		return false
	}

	err3, b := encodeInt(i + 1)
	if err3 {
		return false
	}

	vm.Push(b, false)
	return true
}

func OP_1SUB(vm *VM) bool {
	err1, v := vm.Pop(false)
	if err1 {
		return false
	}

	err2, i := decodeInt(v)
	if err2 {
		return false
	}

	err3, b := encodeInt(i - 1)
	if err3 {
		return false
	}

	vm.Push(b, false)
	return true
}

func OP_2MUL(vm *VM) bool {
	return false // Disabled
}

func OP_2DIV(vm *VM) bool {
	return false // Disabled
}

func OP_NEGATE(vm *VM) bool {
	err1, v := vm.Pop(false)
	if err1 {
		return false
	}

	err2, i := decodeInt(v)
	if err2 {
		return false
	}

	err3, b := encodeInt(i * -1)
	if err3 {
		return false
	}

	vm.Push(b, false)
	return true
}

func OP_ABS(vm *VM) bool {
	err1, v := vm.Pop(false)
	if err1 {
		return false
	}

	err2, i := decodeInt(v)
	if err2 {
		return false
	}

	ii := i >> 63 // http://cavaliercoder.com/blog/optimized-abs-for-int64-in-go.html
	err3, b := encodeInt((i ^ ii) - ii)
	if err3 {
		return false
	}

	vm.Push(b, false)
	return true
}

func OP_NOT(vm *VM) bool {
	err1, v := vm.Pop(false)
	if err1 {
		return false
	}

	err2, i := decodeInt(v)
	if err2 {
		return false
	}

	var b byte = 0x00
	if i == 0 {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_0NOTEQUAL(vm *VM) bool {
	err1, v := vm.Pop(false)
	if err1 {
		return false
	}

	var b byte = 0x01
	if isZero(v) {
		b = 0x00
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_ADD(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	err, b := encodeInt(i1 + i2)
	if err {
		return false
	}

	vm.Push(b, false)
	return true
}

func OP_SUB(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	err, b := encodeInt(i2 - i1)
	if err {
		return false
	}

	vm.Push(b, false)
	return true
}

func OP_MUL(vm *VM) bool {
	return false // Disabled
}

func OP_DIV(vm *VM) bool {
	return false // Disabled
}

func OP_MOD(vm *VM) bool {
	return false // Disabled
}

func OP_LSHIFT(vm *VM) bool {
	return false // Disabled
}

func OP_RSHIFT(vm *VM) bool {
	return false // Disabled
}

func OP_BOOLAND(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	var b byte = 0x00
	if !isZero(v1) && !isZero(v2) {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_BOOLOR(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	var b byte = 0x00
	if !isZero(v1) || !isZero(v2) {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_NUMEQUAL(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	var b byte = 0x00
	if i1 == i2 {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_NUMEQUALVERIFY(vm *VM) bool {
	success := OP_NUMEQUAL(vm)
	return success && OP_VERIFY(vm)
}

func OP_NUMNOTEQUAL(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	var b byte = 0x00
	if i1 != i2 {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_LESSTHAN(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	var b byte = 0x00
	if i2 < i1 {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_GREATERTHAN(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	var b byte = 0x00
	if i2 > i1 {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_LESSTHANOREQUAL(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	var b byte = 0x00
	if i2 <= i1 {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_GREATERTHANOREQUAL(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	var b byte = 0x00
	if i2 >= i1 {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}

func OP_MIN(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	if i2 < i1 {
		vm.Push(v2, false)
	} else {
		vm.Push(v1, false)
	}
	return true
}

func OP_MAX(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	if err1 || err2 {
		return false
	}

	err3, i1 := decodeInt(v1)
	err4, i2 := decodeInt(v2)
	if err3 || err4 {
		return false
	}

	if i2 > i1 {
		vm.Push(v2, false)
	} else {
		vm.Push(v1, false)
	}
	return true
}

func OP_WITHIN(vm *VM) bool {
	err1, v1 := vm.Pop(false)
	err2, v2 := vm.Pop(false)
	err3, v3 := vm.Pop(false)
	if err1 || err2 || err3 {
		return false
	}

	err4, max := decodeInt(v1)
	err5, min := decodeInt(v2)
	err6, x := decodeInt(v3)
	if err4 || err5 || err6 {
		return false
	}

	var b byte = 0x00
	if x >= min && x < max {
		b = 0x01
	}

	vm.Push([]byte{b}, false)
	return true
}
