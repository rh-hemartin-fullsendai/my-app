package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	// Create a temporary index.html for testing
	dir := t.TempDir()
	indexContent := "<!DOCTYPE html><html><body>Test</body></html>"
	err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(indexContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Change to the temp dir so ServeFile finds index.html
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)
	os.Chdir(dir)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	indexHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Test") {
		t.Errorf("expected body to contain 'Test', got %q", body)
	}
}

func TestCounterHandler(t *testing.T) {
	// Reset counter
	counterMu.Lock()
	counter = 0
	counterMu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/counter", nil)
	w := httptest.NewRecorder()
	counterHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if body != "0" {
		t.Errorf("expected '0', got %q", body)
	}
}

func TestCounterIncrementHandler(t *testing.T) {
	// Reset counter
	counterMu.Lock()
	counter = 0
	counterMu.Unlock()

	// First increment
	req := httptest.NewRequest(http.MethodPost, "/counter/increment", nil)
	w := httptest.NewRecorder()
	counterIncrementHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if body != "1" {
		t.Errorf("expected '1', got %q", body)
	}

	// Second increment
	req = httptest.NewRequest(http.MethodPost, "/counter/increment", nil)
	w = httptest.NewRecorder()
	counterIncrementHandler(w, req)

	body = w.Body.String()
	if body != "2" {
		t.Errorf("expected '2', got %q", body)
	}
}

func TestCounterIncrementHandlerRejectsGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/counter/increment", nil)
	w := httptest.NewRecorder()
	counterIncrementHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", resp.StatusCode)
	}
}
