package widevine

import (
	"encoding/hex"
	"encoding/json"
	"testing"
)

// AES key and IV for the provider "widevine_test".
// Use these test keys for testing or integration tests.
var (
	key = []byte{
		0x1a, 0xe8, 0xcc, 0xd0, 0xe7, 0x98, 0x5c, 0xc0,
		0xb6, 0x20, 0x3a, 0x55, 0x85, 0x5a, 0x10, 0x34,
		0xaf, 0xc2, 0x52, 0x98, 0x0e, 0x97, 0x0c, 0xa9,
		0x0e, 0x52, 0x02, 0x68, 0x9f, 0x94, 0x7a, 0xb9}

	iv = []byte{
		0xd5, 0x8c, 0xe9, 0x54, 0x20, 0x3b, 0x7c, 0x9a,
		0x9a, 0x9d, 0x46, 0x7f, 0x59, 0x83, 0x92, 0x49}
)

func TestPad(t *testing.T) {
	text := "deadbeef"

	textDec, _ := hex.DecodeString(text)
	plaintext := pad([]byte(textDec))

	if len(plaintext) != 16 {
		t.Error()
	}
}

func TestUnpad(t *testing.T) {
	hexText := []byte{222, 173, 190, 239, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	out, _ := unpad(hexText)
	dec := hex.EncodeToString(out)

	if len(out) != 4 && dec != "deadbeef" {
		t.Error()
	}
}

func TestEncrypt(t *testing.T) {
	text := "supersecret"

	c := NewCrypto(key, iv)
	enc := c.encrypt(text)

	if enc != "8E2tY1KlW2830Q7EpsjS+A==" {
		t.Error()
	}
}

func TestGenerateSignature(t *testing.T) {

	c := NewCrypto(key, iv)
	payload := map[string]interface{}{
		"test":   "testing",
		"test2":  "testing2",
		"test3":  "testing3",
		"isTest": true,
	}
	jsonPayload, _ := json.Marshal(payload)
	sig := c.generateSignature(jsonPayload)

	if sig != "ga80QzRuUM+jnPcoR6UWs5TXrTQ2VgeYiu0FoqCNRH4=" {
		t.Error()
	}
}
