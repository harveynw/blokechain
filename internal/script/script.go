package script

import (
	"fmt"
	"reflect"
    "runtime"
	"github.com/harveynw/blokechain/internal/script/ops"
)

// VM implements the bitcoin virtual machine
type VM struct {
	Stack [][]byte
	AltStack [][]byte
	Transaction []byte // Required for operations dependent on the transaction
}

// Script for containing and encoding a script
type Script struct {
	data []byte
}

var operations = map[byte]func(*VM)bool {
	// FLOW CONTROL (Branching handled at execution)
	0x61 : ops.OP_NOP,
	0x69 : ops.OP_VERIFY,
	0x6a : ops.OP_RETURN,

	// STACK
	0x6b : ops.OP_TOALTSTACK,
	0x6c : ops.OP_FROMALTSTACK,
	0x73 : ops.OP_IFDUP,
	0x74 : ops.OP_DEPTH,
	0x75 : ops.OP_DROP,
	0x76 : ops.OP_DUP,
	0x77 : ops.OP_NIP,
	0x78 : ops.OP_OVER,
	0x79 : ops.OP_PICK,
	0x7a : ops.OP_ROLL,
	0x7b : ops.OP_ROT,
	0x7c : ops.OP_SWAP,
	0x7d : ops.OP_TUCK,
	0x6d : ops.OP_2DROP,
	0x6e : ops.OP_2DUP,
	0x6f : ops.OP_3DUP,
	0x70 : ops.OP_2OVER,
	0x71 : ops.OP_2ROT,
	0x72 : ops.OP_2SWAP,

	// CRYPTO
	0xa6 : ops.OP_RIPEMD160,
	0xa7 : ops.OP_SHA1,
	0xa8 : ops.OP_SHA256,
	0xa9 : ops.OP_HASH160,
	0xaa : ops.OP_HASH256,
	0xab : ops.OP_CODESEPERATOR,
	0xac : ops.OP_CHECKSIG,
	0xad : ops.OP_CHECKSIGVERIFY,
	0xae : ops.OP_CHECKMULTISIG, // TODO Not implemented!
	0xaf : ops.OP_CHECKMULTISIGVERIFY,


	0x88 : ops.OP_EQUALVERIFY,

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

// NewScript creates an empty script
func NewScript() *Script {
	return &Script{data: make([]byte, 0)}
}

// NewVM creates a new execution environment
func NewVM(transactionEncoded []byte) *VM {
	return &VM{
		Stack: make([][]byte, 0),
		AltStack: make([][]byte, 0),
		Transaction: transactionEncoded,
	}
}

// AppendOpCode for an opcode
func (src *Script) AppendOpCode(op byte) {
	src.data = append(src.data, op)
}

// AppendData for arbitrary bytes
func (src *Script) AppendData(b []byte) {
	src.data = append(src.data, data.EncodeInt(len(b), 1)...)
	src.data = append(src.data, b...)
}

// Execute the script, returns bool depending whether the top element of the stack is truthy
func (src *Script) Execute(transactionEncoded []byte) bool {
	stack := NewVM(transactionEncoded)
	scriptBytes := src.data
	for {
		if len(scriptBytes) == 0 {
			break
		}

		var isOp bool
		var selected []byte
		isOp, selected, scriptBytes = scanNext(scriptBytes)

		if isOp {
			success := operations[selected[0]](stack)
			if !success {
				fmt.Println("Failed on", retrieveOpName(selected[0]))
				return false
			}
		} else {
			stack.Push(selected, false)
		}
	}

	err, top := stack.Pop(false)
	if err || len(top) == 0 || (len(top) == 1 && top[0] == 0x00) {
		return false
	}
	return true
}

// Encode returns script as bytes
func (src *Script) Encode() []byte {
	return src.data
}

// Concat appends one script on to the end of another
func (src *Script) Concat(a *Script) {
	src.data = append(src.data, a.data...)
}

// DecodeScript recovers script object for execution
func DecodeScript(data []byte) *Script {
	return &Script{data: data}
}

// Print displays script in readable format
func (src *Script) Print() {
	fmt.Println("BEGIN BLOKE SCRIPT")
	scriptBytes := src.data
	line := 1
	for {
		if len(scriptBytes) == 0 {
			break
		}

		var isOp bool
		var selected []byte
		isOp, selected, scriptBytes = scanNext(scriptBytes)

		if isOp {
			fmt.Println(line, retrieveOpName(selected[0]))
		} else {
			fmt.Printf("%v %x \n", line, selected)
		}
		line++
	}
}

func retrieveOpName(op byte) string {
	opFunc := operations[op]
	fullName := runtime.FuncForPC(reflect.ValueOf(opFunc).Pointer()).Name()

	for i := 0; i < len(fullName) - 2; i++ {
		if fullName[i:i+3] == "OP_" {
			return fullName[i:]
		}
	}
	return ""
}

// Find and returns next opcode or data as well as the rest of scriptBytes
func scanNext(scriptBytes []byte) (isOp bool, selected []byte, remainingBytes []byte) {
	if len(scriptBytes) == 0 {
		return false, nil, nil
	}

	first := scriptBytes[0]
	if first < 0x4c && first > 0x00 {
		// Data
		return false, scriptBytes[1:int(first)+1], scriptBytes[int(first)+1:]
	}
	// Opcode
	return true, []byte{first}, scriptBytes[1:]
}