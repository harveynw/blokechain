package data

import (
	"math/big"
	"encoding/json"
)

// MarshalJSON is to allow the json package to marshal the public key type
func (pk *PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		X big.Int `json:"x"`
		Y big.Int `json:"y"`
	}{
		X: pk.p.x,
		Y: pk.p.y,
	})
}

// UnmarshalJSON is to allow the json package to unmarshal the public key type
func (pk *PublicKey) UnmarshalJSON(data []byte) error {
	recovered := &struct{
		X big.Int `json:"x"`
		Y big.Int `json:"y"`
	}{}

	if err := json.Unmarshal(data, &recovered); err != nil {
		return err
	}

	p := point{curve: &secp256k1, x: recovered.X, y: recovered.Y}
	pk.p = p
	return nil
}