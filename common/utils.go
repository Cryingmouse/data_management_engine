package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
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
		return "", fmt.Errorf("failed to decrypt the sensitive informaiton")
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

func SplitToList(fields string) []string {
	if fields == "" {
		return []string{} // 返回空切片
	}

	return strings.Split(fields, ",")
}

func CopyStructList(src, dest interface{}) error {
	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	// 如果 src 或 dest 是指针，则获取其指向的值
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	if destVal.Kind() == reflect.Ptr {
		destVal = destVal.Elem()
	}

	// 检查 src 和 dest 是否为结构体列表或结构体
	if srcVal.Kind() == reflect.Slice && destVal.Kind() == reflect.Slice {
		// src 和 dest 是结构体列表
		srcLen := srcVal.Len()

		// 扩展 dest 的长度
		destType := destVal.Type()
		destSlice := reflect.MakeSlice(destType, srcLen, srcLen)
		reflect.Copy(destSlice, destVal)

		// 复制结构体列表
		for i := 0; i < srcLen; i++ {
			srcStruct := srcVal.Index(i)
			destStruct := destSlice.Index(i)

			// 检查 srcStruct 和 destStruct 是否为结构体
			if srcStruct.Kind() == reflect.Ptr {
				srcStruct = srcStruct.Elem()
			}

			if destStruct.Kind() == reflect.Ptr {
				destStruct = destStruct.Elem()
			}

			if srcStruct.Kind() != reflect.Struct || destStruct.Kind() != reflect.Struct {
				return errors.New("src and dest must contain struct instances")
			}

			// 复制结构体字段
			err := copyStructFields(srcStruct, destStruct)
			if err != nil {
				return err
			}
		}

		// 将扩展后的 dest 赋值回原始 dest 变量
		reflect.ValueOf(dest).Elem().Set(destSlice)
	} else if srcVal.Kind() == reflect.Struct && destVal.Kind() == reflect.Struct {
		// src 和 dest 是单个结构体
		return copyStructFields(srcVal, destVal)
	} else {
		return errors.New("src and dest must be either slice of structs or structs")
	}

	return nil
}

func copyStructFields(src, dest reflect.Value) error {
	destType := dest.Type()

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		destField := dest.FieldByName(field.Name)

		if destField.IsValid() && destField.CanSet() {
			srcField := src.FieldByName(field.Name)

			if srcField.IsValid() && srcField.Type() == destField.Type() {
				destField.Set(srcField)
			}
		}
	}

	return nil
}
