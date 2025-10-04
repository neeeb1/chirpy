package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashing(t *testing.T) {
	cases := []struct {
		input    []string
		expected bool
	}{
		{
			input:    []string{"pa$$word", "pa$$word"},
			expected: true,
		},
		{
			input:    []string{"pa$$word", "incorrectpa$$"},
			expected: false,
		},
	}

	for _, c := range cases {
		hash, err := HashPassword(c.input[0])
		if err != nil {
			t.Errorf("failed to hash password: %v", err)
		}
		actual, err := CheckPasswordHash(c.input[1], hash)
		if err != nil {
			t.Errorf("failed to check hash: %v", err)
		}
		if actual != c.expected {
			t.Errorf("failed to check hash: expected %v, got %v", c.expected, actual)
		}
	}
}

func TestJWT(t *testing.T) {
	type result struct {
		UserID uuid.UUID
		err    error
	}

	type input struct {
		userID       uuid.UUID
		tokenSecret  string
		expiresIn    time.Duration
		signedSecret string
	}

	cases := []struct {
		input    input
		expected result
	}{
		{
			input: input{
				userID:       uuid.MustParse("65a8aa5c-b0b4-4240-a6dd-d464bf0a09e7"),
				tokenSecret:  "myspecialsecret",
				expiresIn:    time.Minute * 10,
				signedSecret: "myspecialsecret",
			},
			expected: result{
				UserID: uuid.MustParse("65a8aa5c-b0b4-4240-a6dd-d464bf0a09e7"),
				err:    nil,
			},
		},
		{
			input: input{
				userID:       uuid.MustParse("bf719333-efb5-47c8-9d85-ffb4b6e73996"),
				tokenSecret:  "myspecialsecret",
				expiresIn:    time.Millisecond * 10,
				signedSecret: "myspecialsecret",
			},
			expected: result{
				UserID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				err:    fmt.Errorf("token has invalid claims: token is expired"),
			},
		},
		{
			input: input{
				userID:       uuid.MustParse("bf719333-efb5-47c8-9d85-ffb4b6e73996"),
				tokenSecret:  "myspecialsecret",
				expiresIn:    time.Minute * 10,
				signedSecret: "wrongsecret!",
			},
			expected: result{
				UserID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				err:    fmt.Errorf("token signature is invalid: signature is invalid"),
			},
		},
	}

	for _, c := range cases {
		token, err := MakeJWT(c.input.userID, c.input.tokenSecret, c.input.expiresIn)
		if err != nil {
			t.Errorf("%v", err)
		}

		time.Sleep(time.Millisecond * 100)

		var actual result
		actual.UserID, actual.err = ValidateJWT(token, c.input.signedSecret)

		if c.expected.err == nil && actual.err != nil {
			t.Errorf("failed: expected %v, got %v", c.expected.err, actual.err)
		}
		if c.expected.err != nil && actual.err == nil {
			t.Errorf("failed: expected %v, got %v", c.expected.err, actual.err)
		}
		if c.expected.err != nil && actual.err != nil {
			if c.expected.err.Error() != actual.err.Error() {
				t.Errorf("failed: expected %v, got %v", c.expected.err.Error(), actual.err.Error())
			}
		}

		if actual.UserID != c.expected.UserID {
			t.Errorf("failed to validate token: expected %v, got %v", c.expected.UserID, actual.UserID)
		}

	}
}
