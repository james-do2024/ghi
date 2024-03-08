package client

import (
	"os"
	"testing"
)

func TestNewRestClient(t *testing.T) {
	// Setup environment variables for testing
	os.Setenv("GH_AUTH_TOKEN", "github_pat_dummyToken12345678901234567890")
	defer os.Unsetenv("GH_AUTH_TOKEN") // Clean up after test

	// Can we create a new client without issue?
	client, err := NewRestClient()
	if err != nil {
		t.Fatalf("NewRestClient() error = %v, wantErr %v", err, false)
	}
	if client == nil {
		t.Fatalf("NewRestClient() returned nil client, want non-nil")
	}
}

// Test isPATValid against different formats
func TestIsPATValid(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  bool
	}{
		{"ValidToken", "github_pat_abcdefghijklmnopqrstuvwxyz123456", true},
		{"InvalidTokenShort", "github_pat_short", false},
		{"InvalidTokenLong", "github_pat_" + makeString('a', 244), false}, // assuming makeString repeats a character n times
		{"InvalidTokenBadPrefix", "badprefix_abcdefghijklmnopqrstuvwxyz123456", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPATValid(tt.token); got != tt.want {
				t.Errorf("isPATValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeString(char rune, length int) string {
	return string(make([]rune, length))
}
