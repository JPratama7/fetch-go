package robot

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		path     string
		expected bool
	}{
		{
			name:     "simple allow",
			content:  "User-agent: *\nAllow: /allowed-path",
			path:     "/allowed-path",
			expected: true,
		},
		{
			name:     "simple disallow",
			content:  "User-agent: *\nDisallow: /disallowed",
			path:     "/disallowed",
			expected: false,
		},
		{
			name:     "no rules",
			content:  "",
			path:     "/any-path",
			expected: true,
		},
		{
			name:     "order matters",
			content:  "User-agent: *\nAllow: /api\nDisallow: /api/secret",
			path:     "/api/secret",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := Parse(tt.content)
			result := rules.IsAllowed(tt.path)
			if result != tt.expected {
				t.Errorf("IsAllowed(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestWildcards(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		path     string
		expected bool
	}{
		{
			name:     "wildcard match",
			content:  "User-agent: *\nDisallow: /fish*",
			path:     "/fish.html",
			expected: false,
		},
		{
			name:     "end match",
			content:  "User-agent: *\nAllow: /fish$",
			path:     "/fish",
			expected: true,
		},
		{
			name:     "end match negative",
			content:  "User-agent: *\nAllow: /fish$",
			path:     "/fish.html",
			expected: true, // Default allow when no rules match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := Parse(tt.content)
			result := rules.IsAllowed(tt.path)
			if result != tt.expected {
				t.Errorf("IsAllowed(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestRulePrecedence(t *testing.T) {
	content := `User-agent: *
Allow: /fish
Disallow: /fish/salmon`

	tests := []struct {
		path     string
		expected bool
	}{
		{"/fish", true},
		{"/fish/salmon", false}, // More specific rule takes precedence
		{"/fish/trout", true},   // Only /fish/salmon is disallowed
	}

	rules := Parse(content)
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := rules.IsAllowed(tt.path)
			if result != tt.expected {
				t.Errorf("IsAllowed(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
