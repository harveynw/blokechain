package chain

type Locktime struct {
	t int32
}

func NewLocktime(t int) Locktime {
	return Locktime{t : int32(t)}
}

func (lt Locktime) Encode() []byte {
	// Little endian 4 bytes
	return reverseBytes(EncodeInt(int64(lt.t), 4))
}

func DecodeLocktime(b []byte) Locktime {
	if len(b) != 4 {
		panic("Expected 4 bytes for locktime")
	}
	return Locktime{
		t: int32(DecodeInt(reverseBytes(b))),
	}
}