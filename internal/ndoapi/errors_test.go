package ndoapi

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorForObjectNotFound(t *testing.T) {
	parseErr := errors.New("invalid JSON")
	tests := []struct {
		name         string
		response     *http.Response
		wantNotFound bool
	}{
		{name: "plain 404", response: &http.Response{StatusCode: http.StatusNotFound}, wantNotFound: true},
		{name: "other status", response: &http.Response{StatusCode: http.StatusForbidden}},
		{name: "no response"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := errorForObjectNotFound("api/v1/templates/missing", test.response, parseErr)
			if errors.Is(err, ErrNotFound) != test.wantNotFound {
				t.Fatalf("unexpected not-found classification: %v", err)
			}
			if !test.wantNotFound && !errors.Is(err, parseErr) {
				t.Fatalf("expected original error, got %v", err)
			}
		})
	}
}
