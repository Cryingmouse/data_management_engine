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

// In checks if a target string is present in a sorted string array.
func In(target string, strArray []string) bool {
	index := sort.SearchStrings(strArray, target)
	return index < len(strArray) && strArray[index] == target
}

// AES constants
const (
	KeySize   = 32 // AES-256
	NonceSize = 12 // GCM Nonce Size
	TagSize   = 16 // GCM Tag Size
)

// Encrypt encrypts the plaintext using AES-GCM with the provided key.
func Encrypt(plaintext, key string) (string, error) {
	if len(key) != KeySize {
		return "", errors.New("invalid key length")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	encryptedMessage := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(encryptedMessage), nil
}

// Decrypt decrypts the encryptedText using AES-GCM with the provided key.
func Decrypt(encryptedText, key string) (string, error) {
	if len(key) != KeySize {
		return "", errors.New("invalid key length")
	}

	encryptedMessage, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(encryptedMessage) < NonceSize+TagSize {
		return "", errors.New("invalid encrypted message")
	}

	nonce := encryptedMessage[:NonceSize]
	ciphertext := encryptedMessage[NonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// SplitToList splits the fields string into a string slice.
func SplitToList(fields string) []string {
	if fields == "" {
		return nil
	}

	return strings.Split(fields, ",")
}

// CopyStructList copies the values from src to dest for a slice of structs or single structs.
func CopyStructList(src, dest interface{}) error {
	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	if destVal.Kind() == reflect.Ptr {
		destVal = destVal.Elem()
	}

	if srcVal.Kind() != reflect.Slice && srcVal.Kind() != reflect.Struct {
		return errors.New("src must be a slice or struct")
	}

	if destVal.Kind() != reflect.Slice && destVal.Kind() != reflect.Struct {
		return errors.New("dest must be a slice or struct")
	}

	if srcVal.Kind() == reflect.Slice && destVal.Kind() == reflect.Slice {
		srcLen := srcVal.Len()

		destType := destVal.Type()
		destSlice := reflect.MakeSlice(destType, srcLen, srcLen)
		reflect.Copy(destSlice, destVal)

		for i := 0; i < srcLen; i++ {
			srcStruct := srcVal.Index(i)
			destStruct := destSlice.Index(i)

			if srcStruct.Kind() == reflect.Ptr {
				srcStruct = srcStruct.Elem()
			}

			if destStruct.Kind() == reflect.Ptr {
				destStruct = destStruct.Elem()
			}

			if srcStruct.Kind() != reflect.Struct || destStruct.Kind() != reflect.Struct {
				return errors.New("src and dest must contain struct instances")
			}

			if err := copyStructFields(srcStruct, destStruct); err != nil {
				return err
			}
		}

		reflect.ValueOf(dest).Elem().Set(destSlice)
	} else if srcVal.Kind() == reflect.Struct && destVal.Kind() == reflect.Struct {
		return copyStructFields(srcVal, destVal)
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

func StructListToMapList(structList interface{}) []map[string]interface{} {
	sliceValue := reflect.ValueOf(structList)
	sliceType := sliceValue.Type()

	if sliceType.Kind() != reflect.Slice {
		panic("structListToMapList: input is not a slice")
	}

	mapList := make([]map[string]interface{}, sliceValue.Len())

	for i := 0; i < sliceValue.Len(); i++ {
		structValue := sliceValue.Index(i)
		mapValue := StructToMap(structValue.Interface())
		mapList[i] = mapValue
	}

	return mapList
}

func parseGormTag(tag string) (column string) {
	tagValues := strings.Split(tag, ";")

	for _, value := range tagValues {
		if strings.HasPrefix(value, "column:") {
			column = strings.TrimPrefix(value, "column:")
		}
	}

	return column
}

func StructToMap(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Struct {
		t := v.Type()

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			zeroValue := reflect.Zero(value.Type())

			if reflect.DeepEqual(value.Interface(), zeroValue.Interface()) {
				continue
			}

			column := parseGormTag(field.Tag.Get("gorm"))
			if column != "" {
				result[column] = value.Interface()
			}
		}
	}

	return result
}

func AddQuotes(s string) string {
	return "\"" + s + "\""
}
