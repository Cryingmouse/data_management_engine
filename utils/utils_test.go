package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyStruct(t *testing.T) {
	type Address struct {
		Street string
		City   string
	}

	type Person struct {
		Name    string
		Age     int
		Address Address
	}

	type Employee struct {
		Person
		EmployeeID int
	}

	type Customer struct {
		Person
		CustomerID int
	}

	// 创建源结构体实例
	src := Employee{
		Person: Person{
			Name: "John",
			Age:  30,
			Address: Address{
				Street: "123 Main St",
				City:   "Cityville",
			},
		},
		EmployeeID: 12345,
	}

	var dest Customer

	err := CopyStruct(src, &dest)
	assert.NoError(t, err, "CopyStruct should not return an error")

	assert.Equal(t, src.Name, dest.Name, "Name should be copied correctly")
	assert.Equal(t, src.Age, dest.Age, "Age should be copied correctly")
	assert.Equal(t, src.Address.Street, dest.Address.Street, "Street should be copied correctly")
	assert.Equal(t, src.Address.City, dest.Address.City, "City should be copied correctly")
	assert.Zero(t, dest.CustomerID, "CustomerID should be initialized to zero")
}
