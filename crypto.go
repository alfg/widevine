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

func generateSignature(payload []byte) string {
	h := sha1.New()
	h.Write([]byte(payload))

	bs := h.Sum(nil)
	hash := fmt.Sprintf("%x", bs)

	// Create signature.
	ciphertext := encrypt(hash)
	return ciphertext
}

func encrypt(text string) string {
	// See: https://golang.org/pkg/crypto/cipher/#NewCBCEncrypter

	// key, _ := hex.DecodeString(KEY)
	// iv, _ := hex.DecodeString(IV)

	textb, _ := hex.DecodeString(text)

	plaintext := pad([]byte(textb))

	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	enc := base64.StdEncoding.EncodeToString([]byte(ciphertext[aes.BlockSize:]))
	return enc
}

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}
	return src[:(length - unpadding)], nil
}
