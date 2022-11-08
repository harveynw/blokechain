package script

import (
	"testing"
)

func TestArithmetic(t *testing.T) {
	script := NewScript()

	script.data = []byte{
		0x51, // OP_1
		0x8b, // OP_1ADD
	}
	script.Print()
	result := script.Execute(nil)
	if result != true {
		t.Errorf("Should return true, got %v \n", result)
	}

	script.data = []byte{
		0x8b, // OP_1ADD, try to add one to empty stack
	}
	script.Print()

	result = script.Execute(nil)
	if result != false {
		t.Errorf("Should return false, got %v \n", result)
	}
}

func TestControlFlow(t *testing.T) {
	script := NewScript()

	// OP_2 OP_1ADD OP_IF OP_1SUB OP_1SUB OP_1SUB OP_ENDIF
	script.data = []byte{0x52, 0x8b, 0x63, 0x8c, 0x8c, 0x8c, 0x68}
	script.Print()

	result := script.Execute(nil)
	if result != false {
		t.Errorf("Should return false, got %v \n", result)
	}
}

func TestNestedControlFlow(t *testing.T) {
	script := NewScript()

	// OP_1 OP_IF { OP_1SUB OP_IF OP_1ADD OP_ELSE OP_1SUB OP_ENDIF } OP_ELSE { OP_1ADD } OP_ENDIF OP_DUP OP_EQUAL
	script.data = []byte{
		0x51,
		0x63,
			0x8c,
			0x63,
				0x8b,
			0x67,
				0x8c,
			0x68,
		0x67,
			0x8b,
		0x68,
		0x76,
		0x87,
	}
	script.Print()

	result := script.Execute(nil)
	if result != true {
		t.Errorf("Should return true, got %v \n", result)
	}
}
