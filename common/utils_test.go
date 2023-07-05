package common

import (
	"reflect"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func Test_Copy_StructList(t *testing.T) {
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
