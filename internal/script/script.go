package script

import (
	"fmt"
	"reflect"
	"runtime"
)

// Script for containing and encoding a script
type Script struct {
	data []byte
}

var operations = map[byte]func(*VM) bool{
	// CONSTANTS
	0x00: OP_0,
	0x4f: OP_1NEGATE,
	0x51: OP_TRUE,
	0x52: OP_N(2),
	0x53: OP_N(3),
	0x54: OP_N(4),
	0x55: OP_N(5),
	0x56: OP_N(6),
	0x57: OP_N(7),
	0x58: OP_N(8),
	0x59: OP_N(9),
	0x5A: OP_N(10),
	0x5B: OP_N(11),
	0x5C: OP_N(12),
	0x5D: OP_N(13),
	0x5E: OP_N(14),
	0x5F: OP_N(15),
	0x60: OP_N(16),

	// FLOW CONTROL (Branching handled at execution)
	0x61: OP_NOP,
	0x69: OP_VERIFY,
	0x6a: OP_RETURN,

	// STACK
	0x6b: OP_TOALTSTACK,
	0x6c: OP_FROMALTSTACK,
	0x73: OP_IFDUP,
	0x74: OP_DEPTH,
	0x75: OP_DROP,
	0x76: OP_DUP,
	0x77: OP_NIP,
	0x78: OP_OVER,
	0x79: OP_PICK,
	0x7a: OP_ROLL,
	0x7b: OP_ROT,
	0x7c: OP_SWAP,
	0x7d: OP_TUCK,
	0x6d: OP_2DROP,
	0x6e: OP_2DUP,
	0x6f: OP_3DUP,
	0x70: OP_2OVER,
	0x71: OP_2ROT,
	0x72: OP_2SWAP,

	// SPLICE,
	0x7e: OP_CAT,
	0x7f: OP_SUBSTR,
	0x80: OP_LEFT,
	0x81: OP_RIGHT,
	0x82: OP_SIZE,

	// BITWISE LOGIC
	0x83: OP_INVERT,
	0x84: OP_AND,
	0x85: OP_OR,
	0x86: OP_XOR,
	0x87: OP_EQUAL,
	0x88: OP_EQUALVERIFY,

	// ARITHMETIC
	0x8b: OP_1ADD,
	0x8c: OP_1SUB,
	0x8d: OP_2MUL,
	0x8e: OP_2DIV,
	0x8f: OP_NEGATE,
	0x90: OP_ABS,
	0x91: OP_NOT,
	0x92: OP_0NOTEQUAL,
	0x93: OP_ADD,
	0x94: OP_SUB,
	0x95: OP_MUL,
	0x96: OP_DIV,
	0x97: OP_MOD,
	0x98: OP_LSHIFT,
	0x99: OP_RSHIFT,
	0x9a: OP_BOOLAND,
	0x9b: OP_BOOLOR,
	0x9c: OP_NUMEQUAL,
	0x9d: OP_NUMEQUALVERIFY,
	0x9e: OP_NUMNOTEQUAL,
	0x9f: OP_LESSTHAN,
	0xa0: OP_GREATERTHAN,
	0xa1: OP_LESSTHANOREQUAL,
	0xa2: OP_GREATERTHANOREQUAL,
	0xa3: OP_MIN,
	0xa4: OP_MAX,
	0xa5: OP_WITHIN,

	// CRYPTO
	0xa6: OP_RIPEMD160,
	0xa7: OP_SHA1,
	0xa8: OP_SHA256,
	0xa9: OP_HASH160,
	0xaa: OP_HASH256,
	0xab: OP_CODESEPERATOR,
	0xac: OP_CHECKSIG,
	0xad: OP_CHECKSIGVERIFY,
	0xae: OP_CHECKMULTISIG, // Needs testing
	0xaf: OP_CHECKMULTISIGVERIFY,
}

// NewScript creates an empty script
func NewScript() *Script {
	return &Script{data: make([]byte, 0)}
}

// AppendOpCode for an opcode
func (src *Script) AppendOpCode(op byte) {
	src.data = append(src.data, op)
}

// AppendData for arbitrary bytes
func (src *Script) AppendData(b []byte) {
	src.data = append(src.data, byte(len(b)))
	src.data = append(src.data, b...)
}

// Execute the script, returns success/failure
func (src *Script) Execute(transactionEncoded []byte) bool {
	vm := NewVM(transactionEncoded)
	scriptBytes := src.data

	// Sequentially execute script
	for {
		if len(scriptBytes) == 0 {
			break
		}

		var parserErr bool
		var isOp bool
		var selected []byte
		parserErr, isOp, selected, scriptBytes = parseNext(scriptBytes, vm)
		if parserErr {
			return false
		}

		if isOp {
			fmt.Printf("Stack %x Op:%x(%v) Script Bytes %x \n", vm.Stack, selected, retrieveOpName(selected[0]), scriptBytes)
		} else {
			fmt.Printf("Stack %x Push:%x Script Bytes %x \n", vm.Stack, selected, scriptBytes)
		}

		if isOp {
			success := operations[selected[0]](vm)
			if !success {
				fmt.Println("Failed on", retrieveOpName(selected[0]))
				return false
			}
		} else {
			vm.Push(selected, false)
		}
	}

	// Final result of script
	err, top := vm.Pop(false)
	fmt.Printf("%x \n", top)
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

	scriptBytes := make([]byte, len(src.data))
	copy(scriptBytes, src.data)

	line := 1
	for {
		if len(scriptBytes) == 0 {
			break
		}

		if scriptBytes[0] == op_if {
			fmt.Println(line, "OP_IF")
			scriptBytes = scriptBytes[1:]
		} else if scriptBytes[0] == op_else {
			fmt.Println(line, "OP_ELSE")
			scriptBytes = scriptBytes[1:]
		} else if scriptBytes[0] == op_endif {
			fmt.Println(line, "OP_ENDIF")
			scriptBytes = scriptBytes[1:]
		} else if scriptBytes[0] == op_notif {
			fmt.Println(line, "OP_NOTIF")
			scriptBytes = scriptBytes[1:]
		} else {
			var isOp bool
			var selected []byte
			_, isOp, selected, scriptBytes = parseNext(scriptBytes, nil)
			if isOp {
				fmt.Println(line, retrieveOpName(selected[0]))
			} else {
				fmt.Printf("%v %x \n", line, selected)
			}
		}

		line++
	}
}

func retrieveOpName(op byte) string {
	if op >= 0x52 && op <= 0x60 {
		return "OP_" + fmt.Sprintf("%v", int(op-0x50))
	}
	if op == op_if {
		return "OP_IF"
	} else if op == op_else {
		return "OP_ELSE"
	} else if op == op_endif {
		return "OP_ENDIF"
	} else if op == op_notif {
		return "OP_NOTIF"
	}

	opFunc := operations[op]
	fullName := runtime.FuncForPC(reflect.ValueOf(opFunc).Pointer()).Name()

	for i := 0; i < len(fullName)-2; i++ {
		if fullName[i:i+3] == "OP_" {
			return fullName[i:]
		}
	}
	return ""
}
