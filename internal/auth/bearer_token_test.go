package auth

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		field_1 string
		field_2 string
		field_3 string
		field_4 string
	}{
		{
			field_1: "Authorization",
			field_2: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
		{
			field_1: "Content-Type",
			field_2: "application/json",
		},
		{
			field_1: "Content-Type",
			field_2: "application/json",
			field_3: "Authorization",
			field_4: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
	}

	for _, c := range cases {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set(c.field_1, c.field_2)
		if len(c.field_3) > 0 {
			r.Header.Add(c.field_3, c.field_4)
		}

		token, err := GetBearerToken(r.Header)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			if strings.Contains(token, "Bearer") {
				t.Errorf("%s\n", "'Bearer' should not be in the string")
			} else {
				fmt.Printf("%s\n", token)
			}
		}
	}
}
