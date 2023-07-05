package common

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string
	Age  int
}

func TestCopyStructList(t *testing.T) {
	// Test struct list copying
	srcList := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
	}

	var destList []Person
	err := CopyStructList(srcList, &destList)
	if err != nil {
		t.Errorf("Error copying struct list: %s", err.Error())
	}

	if !reflect.DeepEqual(srcList, destList) {
		t.Errorf("Copied struct list is not equal to the source")
	}

	// Test individual struct copying
	srcStruct := Person{Name: "Alice", Age: 25}
	var destStruct Person
	err = CopyStructList(srcStruct, &destStruct)
	if err != nil {
		t.Errorf("Error copying individual struct: %s", err.Error())
	}

	if !reflect.DeepEqual(srcStruct, destStruct) {
		t.Errorf("Copied individual struct is not equal to the source")
	}
}

func TestIPValidator(t *testing.T) {
	type IPAddress struct {
		IP string `validate:"required,ipvalidator"`
	}

	validate := validator.New()
	validate.RegisterValidation("ipvalidator", IPValidator)

	// Valid IP address
	validIP := IPAddress{IP: "192.168.0.1"}
	err := validate.Struct(validIP)
	assert.NoError(t, err, "Expected no validation error for valid IP address")

	// Loopback IP address
	loopbackIP := IPAddress{IP: "127.0.0.1"}
	err = validate.Struct(loopbackIP)
	assert.Error(t, err, "Expected validation error for loopback IP address")

	// Invalid IP address
	invalidIP := IPAddress{IP: "invalid-ip"}
	err = validate.Struct(invalidIP)
	assert.Error(t, err, "Expected validation error for invalid IP address")
}
