package script

// VM implements the bitcoin virtual machine
type VM struct {
	Stack [][]byte
	AltStack [][]byte
	Transaction []byte // Required for operations dependent on the transaction
}
