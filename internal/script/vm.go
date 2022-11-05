package script

// VM implements the bitcoin virtual machine
type VM struct {
	Stack [][]byte
	AltStack [][]byte
	Transaction []byte // Required for operations dependent on the transaction
}

// NewVM creates a new execution environment
func NewVM(transactionEncoded []byte) *VM {
	return &VM{
		Stack:       make([][]byte, 0),
		AltStack:    make([][]byte, 0),
		Transaction: transactionEncoded,
	}
}

// Push value on top of stack
func (vm *VM) Push(value []byte, alt bool) {
	if alt {
		vm.AltStack = append(vm.AltStack, value)
	} else {
		vm.Stack = append(vm.Stack, value)
	}
}

// Pop value off top of stack
func (vm *VM) Pop(alt bool) (bool, []byte) {
	stack := &vm.Stack
	if alt {
		stack = &vm.AltStack
	}

	// Check we can remove an item
	size := len(*stack)
	if size == 0 {
		return true, nil
	}

	value := (*stack)[size-1]
	*stack = (*stack)[:size-1]

	return false, value
}
