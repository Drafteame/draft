package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 10000
	keySize    = 32
	ivSize     = 16
)

func Encrypt(plaintext string, passphrase string) (string, error) {
	key := pbkdf2.Key([]byte(passphrase), nil, iterations, keySize+ivSize, sha256.New)

	aesKey := key[:keySize]
	iv := key[keySize:]

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	paddedPlaintext := pkcs7Pad([]byte(plaintext), aes.BlockSize)
	ciphertext := make([]byte, len(paddedPlaintext))
	cbc.CryptBlocks(ciphertext, paddedPlaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encrypted string, passphrase string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(passphrase), nil, iterations, keySize+ivSize, sha256.New)

	aesKey := key[:keySize]
	iv := key[keySize:]

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(data))
	cbc.CryptBlocks(plaintext, data)

	unpaddedPlaintext, err := pkcs7Unpad(plaintext)
	if err != nil {
		return "", err
	}

	return string(unpaddedPlaintext), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("empty data")
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("invalid padding size")
	}

	return data[:length-padding], nil
}
