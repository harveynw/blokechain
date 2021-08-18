package script

import (
	"fmt"
	"reflect"
    "runtime"
	"github.com/harveynw/blokechain/internal/data"
)

// ExecutionStack implements simple stack object
type ExecutionStack struct {
	Stack [][]byte
	TransactionEncoded []byte // Required to verify signature in OP_CHECKSIG
}

// Script for containing and encoding a script
type Script struct {
	script []byte
}

var operations = map[byte]func(*ExecutionStack)bool {
	0x76 : OP_DUP,
	0xa9 : OP_HASH160,
	0x88 : OP_EQUALVERIFY,
	0xac : OP_CHECKSIG,
}

// Push value on top of stack
func (s *ExecutionStack) Push(value []byte) {
	s.Stack = append(s.Stack, value)
}

// Pop value off top of stack
func (s *ExecutionStack) Pop() (bool, []byte) {
	size := len(s.Stack)

	if size == 0 {
		return true, nil
	}

	value := (s.Stack)[size-1]
	s.Stack = (s.Stack)[:size-1]
	return false, value
}

// NewScript creates an empty script
func NewScript() *Script {
	script := make([]byte, 0)
	return &Script{script: script}
}

// NewExecutionStack creates an empty stack
func NewExecutionStack(transactionEncoded []byte) *ExecutionStack {
	stack := make([][]byte, 0)
	return &ExecutionStack{Stack: stack, TransactionEncoded: transactionEncoded}
}

// AppendOpCode for an opcode
func (src *Script) AppendOpCode(op byte) {
	src.script = append(src.script, op)
}

// AppendData for arbitrary bytes
func (src *Script) AppendData(b []byte) {
	src.script = append(src.script, data.EncodeInt(len(b), 1)...)
	src.script = append(src.script, b...)
}

// Execute the script, returns bool depending whether the top element of the stack is truthy
func (src *Script) Execute(transactionEncoded []byte) bool {
	stack := NewExecutionStack(transactionEncoded)
	scriptBytes := src.script
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
			stack.Push(selected)
		}
	}

	err, top := stack.Pop()
	if err || len(top) == 0 || (len(top) == 1 && top[0] == 0x00) {
		return false
	}
	return true
}

// Encode returns script as bytes
func (src *Script) Encode() []byte {
	return src.script
}

// DecodeScript recovers script object for execution
func DecodeScript(data []byte) *Script {
	return &Script{script: data}
}

// Print displays script in readable format
func (src *Script) Print() {
	fmt.Println("BEGIN BLOKE SCRIPT")
	scriptBytes := src.script
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

// Concat appends one script on to the end of another
func (src *Script) Concat(a *Script) {
	src.script = append(src.script, a.script...)
}