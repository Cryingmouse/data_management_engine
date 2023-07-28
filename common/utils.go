package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"reflect"
	"regexp"
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

func DeepCopy(src, dest interface{}) error {
	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	if destVal.Kind() == reflect.Ptr {
		destVal = destVal.Elem()
	}

	if srcVal.Kind() != reflect.Slice && srcVal.Kind() != reflect.Struct && srcVal.Kind() != reflect.Map && srcVal.Kind() != reflect.String {
		return errors.New("src must be a slice, struct, map, or string")
	}

	if destVal.Kind() != reflect.Slice && destVal.Kind() != reflect.Struct && destVal.Kind() != reflect.Map && destVal.Kind() != reflect.String {
		return errors.New("dest must be a slice, struct, map, or string")
	}

	if srcVal.Kind() == reflect.Slice && destVal.Kind() == reflect.Slice {
		srcLen := srcVal.Len()

		if srcVal.Type().Elem() == reflect.TypeOf("") && destVal.Type().Elem() == reflect.TypeOf("") {
			// Handle string slices
			destSlice := make([]string, srcLen)
			for i := 0; i < srcLen; i++ {
				destSlice[i] = srcVal.Index(i).String()
			}
			reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(destSlice))
		} else {
			destType := destVal.Type()
			destSlice := reflect.MakeSlice(destType, srcLen, srcLen)
			reflect.Copy(destSlice, destVal)

			for i := 0; i < srcLen; i++ {
				srcElem := srcVal.Index(i)
				destElem := destSlice.Index(i)

				if srcElem.Kind() == reflect.Ptr {
					srcElem = srcElem.Elem()
				}

				if destElem.Kind() == reflect.Ptr {
					destElem = destElem.Elem()
				}

				if srcElem.Kind() != reflect.Struct && srcElem.Kind() != reflect.Map {
					return errors.New("src elements must be struct, map, or string instances")
				}

				if destElem.Kind() != reflect.Struct && destElem.Kind() != reflect.Map {
					return errors.New("dest elements must be struct, map, or string instances")
				}

				if err := copyValue(srcElem, destElem); err != nil {
					return err
				}
			}

			reflect.ValueOf(dest).Elem().Set(destSlice)
		}
	} else if srcVal.Kind() == reflect.Struct && destVal.Kind() == reflect.Struct {
		return copyValue(srcVal, destVal)
	} else if srcVal.Kind() == reflect.Map && destVal.Kind() == reflect.Map {
		// Handle map deep copy
		destMap := reflect.MakeMap(destVal.Type())
		for _, key := range srcVal.MapKeys() {
			srcValue := srcVal.MapIndex(key)
			destValue := reflect.New(destVal.Type().Elem()).Elem()

			if srcValue.Kind() == reflect.Ptr {
				srcValue = srcValue.Elem()
			}

			if destValue.Kind() == reflect.Ptr {
				destValue = destValue.Elem()
			}

			if srcValue.Kind() != reflect.Struct && srcValue.Kind() != reflect.Map {
				return errors.New("src map values must be struct, map, or string instances")
			}

			if destValue.Kind() != reflect.Struct && destValue.Kind() != reflect.Map {
				return errors.New("dest map values must be struct, map, or string instances")
			}

			if err := copyValue(srcValue, destValue); err != nil {
				return err
			}

			destMap.SetMapIndex(key, destValue)
		}

		reflect.ValueOf(dest).Elem().Set(destMap)
	} else if srcVal.Kind() == reflect.String && destVal.Kind() == reflect.String {
		// Copy string directly (strings are immutable)
		destVal.SetString(srcVal.String())
	} else {
		return errors.New("unsupported combination of src and dest types")
	}

	return nil
}

func copyValue(src, dest reflect.Value) error {
	destType := dest.Type()

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		destField := dest.FieldByName(field.Name)

		if destField.IsValid() && destField.CanSet() {
			srcField := src.FieldByName(field.Name)

			if srcField.IsValid() {
				if srcField.Kind() == reflect.Slice {
					if err := DeepCopy(srcField.Interface(), destField.Addr().Interface()); err != nil {
						return err
					}
				} else {
					destField.Set(srcField)
				}
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

func MaskPassword(data interface{}) interface{} {
	value := reflect.ValueOf(data)

	// Handle single struct instance
	if value.Kind() == reflect.Struct {
		return maskPassword(value)
	}

	if value.Kind() == reflect.String {
		// 定义一个正则表达式，匹配 "password": "..."
		re := regexp.MustCompile(`("password":\s*")([^\\"]+)(")`)

		// 使用 ReplaceAllString 方法，将匹配到的 "password" 值替换为 "********"
		maskedInput := re.ReplaceAllString(data.(string), `$1********$3`)

		return maskedInput
	}

	// Handle list of struct instances
	if value.Kind() == reflect.Slice {
		copySlice := reflect.MakeSlice(reflect.SliceOf(value.Type().Elem()), value.Len(), value.Len())

		for i := 0; i < value.Len(); i++ {
			instanceCopy := maskPassword(value.Index(i))
			copySlice.Index(i).Set(reflect.ValueOf(instanceCopy))
		}

		return copySlice.Interface()
	}

	// Return data as-is if it's not a struct or slice
	return data
}

func maskPassword(instanceValue reflect.Value) interface{} {
	copyValue := reflect.New(instanceValue.Type()).Elem()
	copyValue.Set(instanceValue)

	// Modify the Password field in the copy
	passwordField := copyValue.FieldByName("Password")
	if passwordField.IsValid() && passwordField.Kind() == reflect.String {
		passwordField.SetString("******")
	}

	return copyValue.Interface()
}
