package zeaburcaddyextension_test

import (
	"fmt"
	"testing"

	zeaburcaddyextension "github.com/zeabur/zeabur-caddy-extension"
)

func TestParseRedirects(t *testing.T) {
	content := `
# Example of redirects
/           /home       302
/home       /index      301
/about      /about-us   302  Country=US
# This is a comment
/old-path   /new-path   302  Country=UK Language=en
/invalid    path        302
`

	expectedResults := []zeaburcaddyextension.RedirectRule{
		{SourcePath: "/", TargetPath: "/home", StatusCode: 302, Conditions: ""},
		{SourcePath: "/home", TargetPath: "/index", StatusCode: 301, Conditions: ""},
		{SourcePath: "/about", TargetPath: "/about-us", StatusCode: 302, Conditions: "Country=US"},
		{SourcePath: "/old-path", TargetPath: "/new-path", StatusCode: 302, Conditions: "Country=UK Language=en"},
	}

	rules, err := zeaburcaddyextension.ParseRedirects(content)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	fmt.Printf("Rules: %+v\n", rules)

	if len(rules) != len(expectedResults) {
		t.Errorf("Expected %d rules, got %d", len(expectedResults), len(rules))
	}

	for i, rule := range rules {
		if rule.SourcePath != expectedResults[i].SourcePath ||
			rule.TargetPath != expectedResults[i].TargetPath ||
			rule.StatusCode != expectedResults[i].StatusCode ||
			rule.Conditions != expectedResults[i].Conditions {
			t.Errorf("Expected %+v, got %+v", expectedResults[i], rule)
		}
	}
}
