package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
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

// 定义常量
const (
	KeySize   = 32 // AES-256
	NonceSize = 12 // GCM Nonce Size
	TagSize   = 16 // GCM Tag Size
)

func Encrypt(plaintext, key string) (string, error) {
	// 验证密钥长度
	if len(key) != KeySize {
		return "", errors.New("invalid key length")
	}

	// 生成 AES 密码块
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 创建 GCM 密码器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机 nonce
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密明文
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// 拼接 nonce 和 ciphertext
	encryptedMessage := append(nonce, ciphertext...)

	// Base64 编码
	return base64.StdEncoding.EncodeToString(encryptedMessage), nil
}

func Decrypt(encryptedText, key string) (string, error) {
	// 验证密钥长度
	if len(key) != KeySize {
		return "", errors.New("invalid key length")
	}

	// Base64 解码
	encryptedMessage, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// 生成 AES 密码块
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 创建 GCM 密码器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 验证长度
	if len(encryptedMessage) < NonceSize+TagSize {
		return "", errors.New("invalid encrypted message")
	}

	// 提取 nonce 和 ciphertext
	nonce := encryptedMessage[:NonceSize]
	ciphertext := encryptedMessage[NonceSize:]

	// 解密密文
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
