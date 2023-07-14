package common

import (
	"reflect"
	"testing"
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

func TestEncrypt(t *testing.T) {
	type args struct {
		plaintext string
		key       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test_password",
			args: args{
				plaintext: "Password123",
				key:       "MySecretForMagnascale!!!",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.plaintext, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
