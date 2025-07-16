package e2etest

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRedirects(t *testing.T) {
	_, endpoint := TestCaddyContainer(t)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Get(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusFound {
		t.Fatalf("Expected status code %d, but got %d", http.StatusFound, response.StatusCode)
	}

	if response.Header.Get("Location") != "/home" {
		t.Fatalf("Expected redirect to /home, but got %s", response.Header.Get("Location"))
	}
}

func TestHeader(t *testing.T) {
	_, endpoint := TestCaddyContainer(t)

	response, err := http.Get(endpoint + "/test.html")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, response.StatusCode)
	}

	if contentType := response.Header.Get("X-Caddy-Test-Passed"); contentType != "true" {
		t.Fatalf("Expected X-Caddy-Test-Passed to be %s, but got %s", "true", contentType)
	}
}

func TestUnsafePath(t *testing.T) {
	_, endpoint := TestCaddyContainer(t)

	response, err := http.Get(endpoint + "/vendor/unsafe_path")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status code %d, but got %d", http.StatusNotFound, response.StatusCode)
	}
}

// note: this test also covers SPA mode

func TestMpaNotFound(t *testing.T) {
	_, endpoint := TestCaddyContainer(t)

	response, err := http.Get(endpoint + "/invalid_path")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	// WIP: if there is 404.html, it should return 404

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(body), "404 page not found") {
		t.Fatalf("Expected body to contain %s, but got %s", "404 page not found", string(body))
	}
}

func TestRedirectToExternalUrl(t *testing.T) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	_, endpoint := TestCaddyContainer(t)

	response, err := client.Get(endpoint + "/google")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusFound {
		t.Fatalf("Expected status code %d, but got %d", http.StatusFound, response.StatusCode)
	}

	if response.Header.Get("Location") != "https://google.com" {
		t.Fatalf("Expected redirect to https://google.com, but got %s", response.Header.Get("Location"))
	}
}
