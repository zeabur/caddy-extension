package zeaburcaddyextension

import (
	"bufio"
	"log/slog"
	"strconv"
	"strings"
)

// RedirectRule represents a single redirect rule.
type RedirectRule struct {
	SourcePath string
	TargetPath string
	StatusCode int
	Conditions string
}

// ParseRedirects takes a string containing the contents of a _redirects file
// and returns a slice of RedirectRule structs, ignoring any invalid lines.
func ParseRedirects(content string) ([]RedirectRule, error) {
	var rules []RedirectRule
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") { // Skip empty lines and comments
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			slog.Warn("Skipping invalid redirect rule", "line", line)
			continue
		}

		statusCode, err := strconv.Atoi(parts[2])
		if err != nil {
			slog.Warn("Skipping invalid status code in redirect rule", "line", line)
			continue
		}

		if !strings.HasPrefix(parts[0], "/") {
			slog.Warn("Skipping invalid source path in redirect rule", "line", line)
			continue
		}

		if !strings.HasPrefix(parts[1], "/") {
			slog.Warn("Skipping invalid target path in redirect rule", "line", line)
			continue
		}

		rule := RedirectRule{
			SourcePath: parts[0],
			TargetPath: parts[1],
			StatusCode: statusCode,
		}

		// Store conditions as a single string if present
		if len(parts) > 3 {
			rule.Conditions = strings.Join(parts[3:], " ")
		}

		rules = append(rules, rule)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}
