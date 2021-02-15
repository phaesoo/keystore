package models

import (
	"encoding/base64"
	"encoding/json"

	"crypto/sha512"

	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

func hash(s string) (string, error) {
	h := sha512.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}

// Payload is information part of JWT(JWS) token
type Payload struct {
	AccessKey string `json:"access_key"`
	Nonce     string `json:"nonce"`
	QueryHash string `json:"query_hash"`
}

// NewPayload creates a new payload with given arguments
func NewPayload(accessKey, nonce, queryString string) (Payload, error) {
	hashed, err := hash(queryString)
	if err != nil {
		return Payload{}, err
	}
	return Payload{
		AccessKey: accessKey,
		Nonce:     nonce,
		QueryHash: hashed,
	}, nil
}

// Validate payload
func (p *Payload) Validate(queryString string) error {
	h, err := hash(queryString)
	if err != nil {
		return err
	}
	if p.QueryHash != h {
		return errors.New("Invalid query hash")
	}
	return nil
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
