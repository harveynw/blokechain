package script

import "fmt"

var op_if byte = 0x63
var op_notif byte = 0x64
var op_else byte = 0x67
var op_endif byte = 0x68

// Find and returns next OP_X or data and remaining scriptBytes, handling control flow implictly
func parseNext(scriptBytes []byte, vm *VM) (err bool, isOpCode bool, selected []byte, remainingBytes []byte) {
	if len(scriptBytes) == 0 {
		return false, false, nil, nil
	}
	first := scriptBytes[0]

	// Control flow encountered, find branch and recurse
	if first == op_if || first == op_notif {
		err, scriptBytes = trimControlFlow(scriptBytes, vm)
		if err {
			return true, false, nil, nil
		}
		return parseNext(scriptBytes, vm)
	}

	// Data or Opcode
	err = false
	isOpCode, selected, remainingBytes = parseStatement(scriptBytes)
	return
}

func parseStatement(scriptBytes []byte) (isOpCode bool, statement []byte, remainingBytes []byte) {
	first := scriptBytes[0]

	// Data
	if first <= 0x4b && first >= 0x01 {
		length := int(first)
		return false, scriptBytes[1 : length+1], scriptBytes[length+1:]
	}

	// Opcode
	return true, []byte{first}, scriptBytes[1:]
}

func trimControlFlow(scriptBytes []byte, vm *VM) (err bool, trimmed []byte) {
	fmt.Printf("Stack %x trimControlFlow called \n", vm.Stack)
	beginsWithOp, op, scriptBytes := parseStatement(scriptBytes)
	if !beginsWithOp || !(op[0] == op_if || op[0] == op_notif) || len(vm.Stack) == 0 {
		fmt.Println("DEBUG trimControlFlow called but script doesn't begin with control flow")
		return false, nil
	}
	fmt.Printf("trim --> Examining %s\n", retrieveOpName(op[0]))

	// Scan forward to distinguish branches
	if_depth, else_depth := 1, 0
	in_false_branch := false
	true_branch, false_branch := make([]byte, 0), make([]byte, 0)
	for {
		var is_op bool
		var data []byte
		is_op, data, scriptBytes = parseStatement(scriptBytes)
		if is_op {
			fmt.Println("trim --> Examining", retrieveOpName(data[0]))
		} else {
			fmt.Println("trim --> Examining", data)
		}

		if is_op {
			op_code := data[0]
			if op_code == op_if || op_code == op_notif {
				if_depth++
			} else if op_code == op_else {
				else_depth++
				// Switch branch
				if if_depth == 1 && else_depth == 1 {
					in_false_branch = true
					continue
				}
			} else if op_code == op_endif {
				if if_depth > 0 {
					if_depth--
				}
				if else_depth > 0 {
					else_depth--
				}
				// Terminate scan
				if if_depth == 0 && else_depth == 0 {
					break
				}
			}
		}

		if in_false_branch {
			false_branch = append(false_branch, data...)
		} else {
			true_branch = append(true_branch, data...)
		}
	}

	// Return correct branch, as well as the remaining script
	topElementTrue := isTruthy(vm.Stack[0])
	fmt.Println("Top element is", topElementTrue)
	if (op[0] == op_if && topElementTrue) || (op[0] == op_notif && !topElementTrue) {
		// First branch
		fmt.Printf("trim -> Return first branch %x \n", true_branch)
		return false, append(true_branch, scriptBytes...)
	} else {
		// Second branch
		fmt.Printf("trim -> Return second branch %x \n", false_branch)
		return false, append(false_branch, scriptBytes...)
	}
}

// if { 1 0
// 	if { 2 0
// 		if { 3 0

// 		} else { 3 1

// 		} 2 0
// 	} 1 0
// } else { 1 1
// 	if { 2 1
// 		if { 3 1

// 		} else { 3 2

// 		} 2 1
// 	} else { 2 2

// 	} 1 1 
// } 0 0

// if { 1 0
// 	if { 2 0

// 	} else { 2 1

// 	} 1 0
// } 0 0