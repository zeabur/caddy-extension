package zeaburcaddyextension

import (
	"bufio"
	"fmt"
	"strings"
)

// HeaderConfig stores the headers for a specific path
type HeaderConfig struct {
	Path    string
	Headers map[string]string
}

// ParseHeaderConfig parses the header configuration from a given string
func ParseHeaderConfig(config string) ([]HeaderConfig, error) {
	var configs []HeaderConfig
	scanner := bufio.NewScanner(strings.NewReader(config))
	var currentPath string
	headers := make(map[string]string)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Detect path line
		if strings.HasPrefix(line, "/") || strings.HasPrefix(line, "http") {
			// If we have collected headers for the previous path, save them
			if currentPath != "" {
				configs = append(configs, HeaderConfig{
					Path:    currentPath,
					Headers: headers,
				})
				headers = make(map[string]string)
			}
			currentPath = line
		} else {
			// Parse header line
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid header line: %s", line)
			}
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// Add the last collected headers
	if currentPath != "" {
		configs = append(configs, HeaderConfig{
			Path:    currentPath,
			Headers: headers,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}
