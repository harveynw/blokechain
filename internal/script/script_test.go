package script

import "testing"

func TestP2PKH(t *testing.T) {
    got := Multiply(2, 3)
    want := 6

    if want != got {
        t.Errorf("Expected '%d', but got '%d'", want, got)
    }
}