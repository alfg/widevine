package widevine

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

// Crypto struct.
type Crypto struct {
	Key []byte
	IV  []byte
}

// NewCrypto creates a Crypto instance with key and iv.
func NewCrypto(key []byte, iv []byte) *Crypto {
	c := &Crypto{
		Key: key,
		IV:  iv,
	}
	return c
}

func (c *Crypto) generateSignature(payload []byte) string {
	h := sha1.New()
	h.Write([]byte(payload))

	bs := h.Sum(nil)
	hash := fmt.Sprintf("%x", bs)

	// Create signature.
	ciphertext := c.encrypt(hash)
	return ciphertext
}

func (c *Crypto) encrypt(text string) string {
	// See: https://golang.org/pkg/crypto/cipher/#NewCBCEncrypter

	textDec, _ := hex.DecodeString(text)

	plaintext := pad([]byte(textDec))

	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(c.Key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	mode := cipher.NewCBCEncrypter(block, c.IV)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	enc := base64.StdEncoding.EncodeToString([]byte(ciphertext[aes.BlockSize:]))
	return enc
}

// Pads src to a multiple of aes.Blocksize (16) using PKCS #7 standard block padding.
// See http://tools.ietf.org/html/rfc5652#section-6.3.
func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Removes PKCS #7 standard block padding from src.
// See http://tools.ietf.org/html/rfc5652#section-6.3.
// This function is the inverse of pad.
// If the padding is not well-formed, unpad returns nil.
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}
	return src[:(length - unpadding)], nil
}
