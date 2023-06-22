package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"sort"
)

func In(target string, str_array []string) bool {
	sort.Strings(str_array)

	index := sort.SearchStrings(str_array, target)
	//index的取值：[0,len(str_array)]
	//需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

func Encrypt(plaintext, key string) (string, error) {
	// Generate a new AES cipher block using the provided key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Create a new Galois/Counter Mode (GCM) cipher using the block cipher
	// GCM provides authenticated encryption and is generally recommended
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a random nonce (IV)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext using the GCM cipher
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Concatenate the nonce and ciphertext to create the final encrypted message
	encryptedMessage := append(nonce, ciphertext...)

	// Encode the encrypted message in base64 for convenient storage or transmission
	return base64.StdEncoding.EncodeToString(encryptedMessage), nil
}

func Decrypt(encrypted_text, key string) (string, error) {
	// Decode the encrypted message from base64
	encryptedMessage, err := base64.StdEncoding.DecodeString(encrypted_text)
	if err != nil {
		return "", err
	}

	// Generate a new AES cipher block using the provided key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()

	if len(encryptedMessage) < nonceSize {
		return "", fmt.Errorf("failed to decrypt the sensitive informaiton.")
	}

	nonce := encryptedMessage[:nonceSize]
	ciphertext := encryptedMessage[nonceSize:]

	// Decrypt the ciphertext using the GCM cipher
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
