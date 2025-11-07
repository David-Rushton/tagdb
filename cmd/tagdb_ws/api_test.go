package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/tagdb"
)

// TODO: Add more tests.
func Test_List(t *testing.T) {
	// Arrange
	configTestEnvironment(t)
	request := httptest.NewRequest("GET", "/list?tags=tag1,tag2", nil)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(List)

	// Act
	handler.ServeHTTP(response, request)

	// Assert
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want %v", status, http.StatusOK)
	}
}

func configTestEnvironment(t *testing.T) {
	testDir := t.TempDir()

	t.Setenv("TAGDB_PORT", "11981")
	t.Setenv("TAGDB_STORAGE_ROOT", testDir)
	t.Setenv("TAGDB_STORAGE_WAL_ROLL_AFTER_BYTES", "1024")
	t.Setenv("TAGDB_STORAGE_BACKGROUND_TASK_INTERVAL_MS", "0")

	config := getConfig()

	startDatabase(config, t.Context())

	// Ensure database is properly closed before test cleanup
	t.Cleanup(func() {
		// Requires explicit stop as t.Context will not trigger <-ctx.Done()
		tagdb.Stop()
	})
}
