package models

import (
	"encoding/json"

	"github.com/square/go-jose"
)

// Payload is information part of JWT token
type Payload struct {
	AccessKey string `json:"access_key"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

// Encrypt payload with given key
func (p *Payload) Encrypt(key string) (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: []byte(key)}, nil)
	if err != nil {
		return "", err
	}
	sig, err := signer.Sign(b)
	if err != nil {
		return "", err
	}
	return sig.FullSerialize(), nil
}
