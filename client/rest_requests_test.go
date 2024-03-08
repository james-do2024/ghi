package client

import (
	"strings"
	"testing"
)

// Deliberately avoiding mocking up the GH API for tests - too much work for this situation

func TestNavigateUp(t *testing.T) {
	tests := []struct {
		initialPath string
		expected    string
	}{
		{"path/to/dir", "path/to"},
		{"path", ""},
		{"", ""},
	}

	for _, tt := range tests {
		req := RestRequest{CurrentPath: tt.initialPath}
		req.NavigateUp()
		if req.CurrentPath != tt.expected {
			t.Errorf("After NavigateUp(), expected '%s', got '%s'", tt.expected, req.CurrentPath)
		}
	}
}

func TestNavigateRoot(t *testing.T) {
	tests := []struct {
		initialPath string
	}{
		{"path/to/dir"},
		{"another/path"},
		{""},
	}

	for _, tt := range tests {
		req := RestRequest{CurrentPath: tt.initialPath}
		req.NavigateRoot()
		if req.CurrentPath != "" {
			t.Errorf("After NavigateRoot(), expected root path '', got '%s'", req.CurrentPath)
		}
	}
}

func TestNavigateIndex(t *testing.T) {
	dirMap := map[int]string{1: "dir1", 2: "dir2", 3: "mordor"}
	tests := []struct {
		index    int
		expected string
	}{
		{1, "dir1"},
		{2, "dir2"},
		{3, "mordor"},
	}

	for _, tt := range tests {
		req := RestRequest{CurrentPath: "path/to"}
		req.NavigateIndex(tt.index, dirMap)
		if tt.expected == "" && req.CurrentPath != "path/to" {
			// Expect no change on invalid index
			t.Errorf("After NavigateIndex() with invalid index, expected '%s', got '%s'", "path/to", req.CurrentPath)
		} else if tt.expected != "" && !strings.HasSuffix(req.CurrentPath, tt.expected) {
			// Directory name should be appended
			t.Errorf("After NavigateIndex() to %d, expected suffix '%s', got '%s'", tt.index, tt.expected, req.CurrentPath)
		}
	}
}
