package cryptography

import (
	"crypto/sha1"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

func RIPEMD160(b []byte) []byte {
	md := ripemd160.New()
	md.Write(b)
	return md.Sum(nil)
}

func SHA256(b []byte) []byte {
	result := sha256.Sum256(b)
	return result[:]
}

func SHA1(b []byte) []byte {
	result := sha1.Sum(b)
	return result[:]
}

// Hash160 produces a 160-bit digest
func Hash160(b []byte) []byte {
	return RIPEMD160(SHA256(b))
}

// Hash256 produces a 256-bit digest
func Hash256(b []byte) []byte {
	return SHA256(SHA256(b))
}
