package common

import (
	"reflect"
	"testing"
)

type Person struct {
	Name       string
	Age        int
	Hobbies    []string
	Attributes map[string]string
}

type Person2 struct {
	Name       string
	Age        int
	Hobbies    []string
	Attributes map[string]string
}

func TestCopySlice(t *testing.T) {
	// Sample data: a list of Person structs, each with a string slice and map property
	srcSlice := []Person{
		{Name: "Alice", Age: 30, Hobbies: []string{"Reading", "Painting"}, Attributes: map[string]string{"Height": "170cm", "Weight": "60kg"}},
		{Name: "Bob", Age: 25, Hobbies: []string{"Gardening", "Cooking"}, Attributes: map[string]string{"Height": "180cm", "Weight": "70kg"}},
		{Name: "Charlie", Age: 35, Hobbies: []string{"Fishing", "Hiking"}, Attributes: map[string]string{"Height": "175cm", "Weight": "75kg"}},
	}

	// Copy the slice
	var destSlice []Person2
	err := DeepCopy(srcSlice, &destSlice)
	if err != nil {
		t.Fatalf("Error copying slice: %s", err)
	}

	destSlice[0].Attributes["Height"] = "190cm"

	// Check if the two slices are equal
	if !reflect.DeepEqual(srcSlice, destSlice) {
		t.Errorf("Copied slice does not match the source slice.")
	}

	// Modify the original source slice and make sure the destination is not affected
	srcSlice[0].Name = "Modified Name"
	srcSlice[0].Hobbies[0] = "Modified Hobby"
	srcSlice[0].Attributes["Height"] = "165cm"

	if srcSlice[0].Name == destSlice[0].Name {
		t.Errorf("Modifying the source slice affected the destination slice: Name field.")
	}

	if srcSlice[0].Hobbies[0] == destSlice[0].Hobbies[0] {
		t.Errorf("Modifying the source slice affected the destination slice: Hobbies field.")
	}

	if srcSlice[0].Attributes["Height"] == destSlice[0].Attributes["Height"] {
		t.Errorf("Modifying the source slice affected the destination slice: Attributes map.")
	}
}

func Test_Encrypt_Decrypt(t *testing.T) {
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
				key:       "0123456789ABCDEF0123456789ABCDEF",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt(tt.args.plaintext, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			plaintext, err := Decrypt(encrypted, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if plaintext != tt.args.plaintext {
				t.Errorf("Encrypt()/Decrypt() = %v, want %v", plaintext, tt.args.plaintext)
			}

		})
	}
}
