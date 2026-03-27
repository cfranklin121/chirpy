package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJwt(t *testing.T) {
	cases := []struct {
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
	}{
		{
			userID:      uuid.MustParse("5927b1ae-55d6-4d53-baa7-bdf3b601a95e"),
			tokenSecret: "abcde",
			expiresIn:   time.Hour * 2,
		},
		{
			userID:      uuid.MustParse("f0f87ec2-a8b5-48cc-b66a-a85ce7c7b862"),
			tokenSecret: "lkdsaglkes",
			expiresIn:   time.Hour * 1,
		},
		{
			userID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			tokenSecret: "lkickeps354090",
			expiresIn:   -time.Hour,
		},
	}

	fmt.Println("****Test valid signature****")
	for _, c := range cases {
		result, err := MakeJWT(c.userID, c.tokenSecret, c.expiresIn)
		if err != nil {
			t.Errorf("MakeJWT error: %s", err)
		} else {
			fmt.Printf("Return JWT: %s\n", result)
			fmt.Print("Pass")
		}

		fmt.Println("----------------------------")

		validate, err := ValidateJWT(result, c.tokenSecret)
		if err != nil {
			fmt.Printf("Validate error: %s\n", err)
			fmt.Print("Pass")
		} else {
			fmt.Printf("Validate: %s\n", validate)
			fmt.Print("Pass")
		}

		fmt.Println()
		fmt.Println("=============================")
		fmt.Println()
	}

	fmt.Println("****Test invalid signature****")

	result, err := MakeJWT(cases[0].userID, cases[0].tokenSecret, cases[0].expiresIn)
	if err != nil {
		t.Errorf("MakeJWT error: %s\n", err)
	}

	validate, err := ValidateJWT(result, cases[1].tokenSecret)
	if err != nil {
		fmt.Printf("Invalid secret: %s\n", err)
		fmt.Print("Pass")
	} else {
		t.Errorf("FAIL: should not validate: %s\n", validate)
	}
}
