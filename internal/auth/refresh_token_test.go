package auth

import (
	"log"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	for i := 0; i < 4; i++ {
		returned_string := MakeRefreshToken()
		log.Printf("%s\n", returned_string)
	}
}
