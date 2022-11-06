package script


// Find and returns next opcode or data as well as the rest of scriptBytes
func scanNext(scriptBytes []byte) (isOp bool, selected []byte, remainingBytes []byte) {
	if len(scriptBytes) == 0 {
		return false, nil, nil
	}

	first := scriptBytes[0]
	if first <= 0x4b && first >= 0x01 {
		// Data
		return false, scriptBytes[1 : int(first)+1], scriptBytes[int(first)+1:]
	}
	// Opcode
	return true, []byte{first}, scriptBytes[1:]
}