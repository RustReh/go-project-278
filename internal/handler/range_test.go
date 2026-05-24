package handler

import (
	"testing"

	"github.com/RustReh/go-project-278/internal/apperr"
)

func TestParseRangeQuery(t *testing.T) {
	tests := []struct {
		raw       string
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{"[0,10]", 0, 10, false},
		{"[5, 10]", 5, 10, false},
		{"[0,0]", 0, 0, false},
		{"", 0, 0, true},
		{"[0,9]", 0, 9, false},
		{"[10,5]", 0, 0, true},
		{"invalid", 0, 0, true},
	}

	for _, tt := range tests {
		start, end, err := parseRangeQuery(tt.raw)
		if tt.wantErr {
			if err == nil {
				t.Fatalf("range %q: expected error", tt.raw)
			}
			if _, ok := apperr.AsAppError(err); !ok {
				t.Fatalf("range %q: expected AppError", tt.raw)
			}
			continue
		}
		if err != nil {
			t.Fatalf("range %q: %v", tt.raw, err)
		}
		if start != tt.wantStart || end != tt.wantEnd {
			t.Fatalf("range %q: got [%d,%d), want [%d,%d)", tt.raw, start, end, tt.wantStart, tt.wantEnd)
		}
	}
}

func TestContentRangeHeader(t *testing.T) {
	got := contentRangeHeader(0, 10, 42)
	want := "links 0-10/42"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
