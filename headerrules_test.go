package zeaburcaddyextension_test

import (
	"reflect"
	"testing"

	zeaburcaddyextension "github.com/zeabur/caddy-extension"
)

// TestParseConfig tests the ParseConfig function
func TestParseConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []zeaburcaddyextension.HeaderConfig
		wantErr  bool
	}{
		{
			name: "Basic test",
			input: `
/templates/index.html
  X-Frame-Options: DENY
  X-XSS-Protection: 1; mode=block
/templates/index2.html
  X-Frame-Options: SAMEORIGIN`,
			expected: []zeaburcaddyextension.HeaderConfig{
				{
					Path: "/templates/index.html",
					Headers: map[string]string{
						"X-Frame-Options":  "DENY",
						"X-XSS-Protection": "1; mode=block",
					},
				},
				{
					Path: "/templates/index2.html",
					Headers: map[string]string{
						"X-Frame-Options": "SAMEORIGIN",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test with comments and empty lines",
			input: `
# a path:
/templates/index.html
  # headers for that path:
  X-Frame-Options: DENY

# another path:
/templates/index2.html
  X-Frame-Options: SAMEORIGIN

# This is a comment
/secure/page
  X-Frame-Options: DENY
  X-Content-Type-Options: nosniff
  Referrer-Policy: no-referrer`,
			expected: []zeaburcaddyextension.HeaderConfig{
				{
					Path: "/templates/index.html",
					Headers: map[string]string{
						"X-Frame-Options": "DENY",
					},
				},
				{
					Path: "/templates/index2.html",
					Headers: map[string]string{
						"X-Frame-Options": "SAMEORIGIN",
					},
				},
				{
					Path: "/secure/page",
					Headers: map[string]string{
						"X-Frame-Options":        "DENY",
						"X-Content-Type-Options": "nosniff",
						"Referrer-Policy":        "no-referrer",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid header line",
			input: `
/templates/index.html
  X-Frame-Options DENY`, // Missing colon
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := zeaburcaddyextension.ParseHeaderConfig(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseConfig() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
