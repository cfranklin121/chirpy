package auth

import (
	"fmt"
	"testing"
)

func TestAuth(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "password",
			expected: "",
		},
		{
			input:    "This_is_a_better_Password",
			expected: "",
		},
		{
			input:    "123456",
			expected: "",
		},
		{
			input:    "khKHDG)*35#$^",
			expected: "",
		},
	}

	hashes := []string{}
	for _, c := range cases {
		fmt.Println("Input:", c.input)
		result, err := HashPassword(c.input)
		fmt.Println("Result:", result)
		hashes = append(hashes, result)
		if err != nil {
			t.Errorf("%s", err)
		}
	}

	fmt.Println("--------------------------------------")

	for i, c := range cases {
		match, err := CheckPasswordHash(c.input, hashes[i])
		if match {
			fmt.Println("Match")
		}
		if err != nil {
			t.Errorf("%s", err)
		}
	}

}
