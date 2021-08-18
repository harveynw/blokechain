package script

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