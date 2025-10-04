package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestHealthHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerHealth)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expected := `status: OK!`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
