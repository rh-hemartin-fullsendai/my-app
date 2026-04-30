package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestDB initializes a temporary SQLite database for testing.
func setupTestDB(t *testing.T) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	if err := initDB(dbPath); err != nil {
		t.Fatalf("failed to init test DB: %v", err)
	}
	t.Cleanup(func() {
		closeDB()
	})
}

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
	setupTestDB(t)

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
	setupTestDB(t)

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
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/counter/increment", nil)
	w := httptest.NewRecorder()
	counterIncrementHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestCounterPersistence(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "persist.db")

	// First session: initialize and increment
	if err := initDB(dbPath); err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}
	if _, err := incrementCounter(); err != nil {
		t.Fatalf("failed to increment: %v", err)
	}
	if _, err := incrementCounter(); err != nil {
		t.Fatalf("failed to increment: %v", err)
	}
	closeDB()

	// Second session: reopen and verify value persisted
	if err := initDB(dbPath); err != nil {
		t.Fatalf("failed to reopen DB: %v", err)
	}
	defer closeDB()

	value, err := getCounter()
	if err != nil {
		t.Fatalf("failed to get counter: %v", err)
	}
	if value != 2 {
		t.Errorf("expected counter to persist as 2, got %d", value)
	}
}
