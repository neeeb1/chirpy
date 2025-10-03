package api

import (
	"fmt"
	"testing"
)

func TestValidateChirp(t *testing.T) {
	cases := []struct {
		input    chirp
		expected struct {
			chirp chirp
			err   error
		}
	}{
		{
			input: chirp{
				Body:   "This is a long chirp, in fact it's way too long to be considered a chirp! It shouldn't be validated at all! Chirps this long should not be allowed on the platform.",
				UserID: "00000000-0000-0000-0000-000000000000",
			},
			expected: struct {
				chirp chirp
				err   error
			}{
				chirp: chirp{},
				err:   fmt.Errorf("chirp is too long"),
			},
		},
		{
			input: chirp{
				Body:   "This is a fornax long chirp, in fact it's way too long to be considered a chirp! It shouldn't be validated at all! Chirps this long should not be allowed on the platform.",
				UserID: "00000000-0000-0000-0000-000000000000",
			},
			expected: struct {
				chirp chirp
				err   error
			}{
				chirp: chirp{},
				err:   fmt.Errorf("chirp is too long"),
			},
		},
		{
			input: chirp{
				Body:   "You sharbert you made me kerfuffle my fornax!",
				UserID: "00000000-0000-0000-0000-000000000000",
			},
			expected: struct {
				chirp chirp
				err   error
			}{
				chirp: chirp{
					Body:   "You **** you made me **** my fornax!",
					UserID: "00000000-0000-0000-0000-000000000000",
				},
				err: nil,
			},
		},
	}
	type result struct {
		chirp chirp
		err   error
	}

	for _, c := range cases {
		var actual result
		actual.chirp, actual.err = validateChirp(c.input)
		if actual.chirp.Body != c.expected.chirp.Body {
			t.Errorf("failed to validate chirp: got %v, wanted %v\n", actual, c.expected)
		}
	}
}
