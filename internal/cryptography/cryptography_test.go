package cryptography

import (
	"fmt"
	"testing"
)

func TestCurveGenerator(t *testing.T) {
	if !g.isOnCurve() {
		t.Errorf("%v is not on %v \n", g, secp256k1)
	}
}

func TestRandPublicKeyOnCurve(t *testing.T) {
	secretKey := gen.randomSecretKey()
	pk := gen.publicKeyFromSecretKey(secretKey)

	fmt.Println("Secret Key:", secretKey)
	fmt.Println("Public Key X:", pk.p.x)
	fmt.Println("Public Key Y:", pk.p.y)

	if !pk.p.isOnCurve() {
		t.Errorf("Rand public key %v is not on secp256k1 \n", pk)
	}
}

func TestSignature(t *testing.T) {
	secretKey := gen.randomSecretKey()
	pk := gen.publicKeyFromSecretKey(secretKey)

	message := []byte("I'm afraid there is no money")
	sig := SignMessage(secretKey, message)

	if !sig.VerifySignature(pk, message) {
		t.Errorf("Public key %v should match signature \n %v \n of messsage \n %v \n", pk, sig, string(message))
	}

	sigRecovered, err := DecodeSignature(sig.Encode())
	if err != nil {
		t.Error("Failed on signature decode")
	}
	if sigRecovered.r.Cmp(sig.r) != 0 || sigRecovered.s.Cmp(sig.s) != 0 {
		t.Error("DER Encoding and then decoding did not give back the same signature")
	}
}