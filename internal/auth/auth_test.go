package auth

import "testing"

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
