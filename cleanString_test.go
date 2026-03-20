package main

import (
	"fmt"
	"testing"
)

func TestCleanString(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			input:    "This is a major fornax",
			expected: "This is a major ****",
		},
		{
			input:    "This is a kerfuffle opinion I need to share with the world",
			expected: "This is a **** opinion I need to share with the world",
		},
		{
			input:    "Oh sharbert thats hot.",
			expected: "Oh **** thats hot.",
		},
		{
			input:    "OH FORNAX",
			expected: "OH ****",
		},
		{
			input:    "kerfuffle!",
			expected: "kerfuffle!",
		},
	}

	for _, c := range cases {
		fmt.Println("Case:", c.input)
		actual := cleanString(c.input)
		fmt.Printf("Expected: %s\n", c.expected)
		fmt.Printf("Actual  : %s\n", actual)

		if c.expected == actual {
			fmt.Println("Pass")
		} else {
			fmt.Println("Fail")
			t.Errorf("Fail")
		}
		fmt.Println("=================")
		fmt.Println()
	}
}
