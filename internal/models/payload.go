package models

type Payload struct {
	AccessKey string `json:"access_key"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}
