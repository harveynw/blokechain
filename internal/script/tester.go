package script

import "fmt"

func TestScript() {
	fmt.Println("Testing script")
	script := NewScript()

	// // OP_2 OP_1ADD OP_IF OP_1SUB OP_1SUB OP_1SUB OP_ENDIF OP_1ADD
	// script.data = []byte{0x52, 0x8b, 0x63, 0x8c, 0x8c, 0x8c, 0x68, 0x8b}

	// // OP_0 OP_IF { OP_1SUB OP_1SUB OP_1SUB} OP_ELSE { OP_1ADD } OP_ENDIF OP_1ADD
	// script.data = []byte{
	// 	0x00,
	// 	0x63,
	// 		0x8c, 0x8c, 0x8c,
	// 	0x67,
	// 		0x8b,
	// 	0x68,
	// 	0x8b,
	// }

	// // OP_0 OP_NOTIF { OP_1SUB OP_1SUB OP_1SUB} OP_ELSE { OP_1ADD } OP_ENDIF OP_1ADD
	// script.data = []byte{
	// 	0x00,
	// 	0x64,
	// 		0x8c, 0x8c, 0x8c,
	// 	0x67,
	// 		0x8b,
	// 	0x68,
	// 	0x8b,
	// }

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

	// // OP_1ADD
	// script.data = []byte{
	// 	0x8b,
	// }

	// script.Print()
	result := script.Execute(nil)
	fmt.Println(result)
}