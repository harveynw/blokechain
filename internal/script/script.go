package script

import (
	"fmt"
	"bytes"
	"encoding/hex"
	"github.com/harveynw/blokechain/internal/data"
)

// ExecutionStack implements simple stack object
type ExecutionStack struct {
	stack [][]byte
	transactionEncoded []byte // Required to verify signature in OP_CHECKSIG
}

// Script for containing and encoding a script
type Script struct {
	script []byte
}

// Push value on top of stack
func (s *ExecutionStack) Push(value []byte) {
	s.stack = append(s.stack, value)
}

// Pop value off top of stack
func (s *ExecutionStack) Pop() (bool, []byte) {
	size := len(s.stack)

	if size == 0 {
		return true, nil
	}

	value := (s.stack)[size-1]
	s.stack = (s.stack)[:size-1]
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
	return &ExecutionStack{stack: stack, transactionEncoded: transactionEncoded}
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
			success := ops[selected[0]](stack)
			if !success {
				fmt.Println("Failed on", opCodeNames[selected[0]])
				return false
			}
		} else {
			stack.Push(selected)
		}

		fmt.Println("STACK:", stack.stack)
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
			fmt.Println(line, opCodeNames[selected[0]])
		} else {
			fmt.Printf("%v %x \n", line, selected)
		}
		line++
	}
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

var opCodeNames = map[byte]string {
	0x76 : "OP_DUP",
	0xa9 : "OP_HASH160",
	0x88 : "OP_EQUALVERIFY",
	0xac : "OP_CHECKSIG",
}

var ops = map[byte]func(*ExecutionStack)bool {
	0x76 : func(s *ExecutionStack) bool {
		// OP_DUP
		err, value := s.Pop()
		if err {
			return false
		}
		s.Push(value)
		s.Push(value)
		return true
	},
	0xa9 : func(s *ExecutionStack) bool {
		// OP_HASH160
		err, value := s.Pop()
		if err {
			return false
		}
		s.Push(data.DoubleHash(value, true))
		return true
	},
	0x88 : func(s *ExecutionStack) bool {
		// OP_EQUALVERIFY
		err1, val1 := s.Pop()
		err2, val2 := s.Pop()
		if err1 || err2 {
			return false
		}
		return bytes.Compare(val1, val2) == 0
	},
	0xac : func(s *ExecutionStack) bool {
		// OP_CHECKSIG
		err1, pubKeyBytes := s.Pop()
		err2, sigBytes := s.Pop()
		if err1 || err2 {
			return false
		}

		sig, err3 := data.DecodeSignature(sigBytes)
		if err3 != nil {
			return false
		}
		pubKey, err4 := data.DecodePublicKey(pubKeyBytes)
		if err4 != nil {
			return false
		}

		if sig.VerifySignature(pubKey, s.transactionEncoded) {
			s.Push([]byte{0x01}) // Truthy
		} else {
			s.Push([]byte{}) // False
		}

		return true
	},
}

// Concat appends one script on to the end of another
func (src *Script) Concat(a *Script) {
	src.script = append(src.script, a.script...)
}

// P2PKH (Pay to Public Key Hash) generates the boilerplate fund locking script
func P2PKH(address []byte) *Script {
	script := NewScript()
	script.AppendOpCode(0x76)
	script.AppendOpCode(0xa9)
	script.AppendData(address)
	script.AppendOpCode(0x88)
	script.AppendOpCode(0xac)
	return script
}

// Test is a utility function for testing this package
func Test() {
	script := NewScript()

	sig := []byte("signaturederp")
	script.AppendData(sig)
	pubkey, _ := hex.DecodeString("932903290494348")
	script.AppendData(pubkey)

	script.AppendOpCode(0x76)
	script.AppendOpCode(0xa9)
	addr, _ := hex.DecodeString("1e5518889b0d3554fe7cd3378ade632aff3069d8")
	script.AppendData(addr)
	script.AppendOpCode(0x88)
	script.AppendOpCode(0xac)

	script.Print()

	fmt.Println(script.Execute([]byte("encoded transaction")))
}