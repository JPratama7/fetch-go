package robot

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Robot struct {
	agents []string
}

// RobotsTxtNotFoundError is returned when a robots.txt file is not found (404).
var RobotsTxtNotFoundError = errors.New("robots.txt not found")

// Rule represents a single allow/disallow rule with wildcard support
type Rule struct {
	Pattern string // Original pattern from robots.txt
	Path    string // Normalized path
	Allow   bool
}

// Rules represents a set of rules for a user-agent
type Rules struct {
	UserAgent string
	Rules     []Rule
}

// FromURL fetches and parses a robots.txt file from a given URL.
func FromURL(client *http.Client, parsedURL *url.URL) (*Rules, error) {
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)

	resp, err := client.Get(robotsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, RobotsTxtNotFoundError
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch robots.txt: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return Parse(string(body)), nil
}

// normalizePath ensures paths are comparable and handles wildcards
func normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	// Remove trailing wildcards (Google treats them as insignificant)
	if strings.HasSuffix(path, "*") {
		path = strings.TrimSuffix(path, "*")
	}
	return path
}

// matchesPath checks if a URL path matches the rule pattern
func (r *Rule) matchesPath(urlPath string) bool {
	// Handle exact match with $
	if strings.HasSuffix(r.Pattern, "$") {
		pattern := strings.TrimSuffix(r.Pattern, "$")
		return urlPath == pattern
	}

	// Handle wildcards
	if strings.Contains(r.Pattern, "*") {
		pattern := strings.ReplaceAll(r.Pattern, "*", "")
		return strings.HasPrefix(urlPath, pattern)
	}

	// Default prefix match
	return strings.HasPrefix(urlPath, r.Path)
}

// Parse parses the contents of a robots.txt file with Google's rules
func Parse(content string) *Rules {
	rules := &Rules{
		UserAgent: "*",
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "User-agent:") {
			agent := strings.TrimSpace(strings.TrimPrefix(line, "User-agent:"))
			if agent == "*" {
				rules.UserAgent = agent
			}
		} else if strings.HasPrefix(line, "Allow:") {
			pattern := strings.TrimSpace(strings.TrimPrefix(line, "Allow:"))
			rules.Rules = append(rules.Rules, Rule{
				Pattern: pattern,
				Path:    normalizePath(pattern),
				Allow:   true,
			})
		} else if strings.HasPrefix(line, "Disallow:") {
			pattern := strings.TrimSpace(strings.TrimPrefix(line, "Disallow:"))
			rules.Rules = append(rules.Rules, Rule{
				Pattern: pattern,
				Path:    normalizePath(pattern),
				Allow:   false,
			})
		}
	}

	return rules
}

// IsAllowed checks if a given path is allowed by the rules with proper precedence
func (r *Rules) IsAllowed(urlPath string) bool {
	urlPath = normalizePath(urlPath)
	var lastMatchingRule *Rule
	var longestMatchLength int

	for _, rule := range r.Rules {
		if rule.matchesPath(urlPath) {
			// Prefer longer/more specific matches
			if len(rule.Path) > longestMatchLength {
				lastMatchingRule = &rule
				longestMatchLength = len(rule.Path)
			}
		}
	}

	if lastMatchingRule == nil {
		return true // Default to allowed if no rules match
	}

	return lastMatchingRule.Allow
}
