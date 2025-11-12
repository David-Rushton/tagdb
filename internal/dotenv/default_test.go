package dotenv

import (
	"os"
	"testing"
)

func Test_loadEnvFiles(t *testing.T) {
	// Arrange.

	// Lowest priority: OS env variables.
	os.Setenv("KEY_0", "VALUE_0")
	os.Setenv("KEY_1", "SHOULD_BE_OVERRIDDEN")

	// Medium priority: .env file.
	envFile := `# comment to ignore
KEY_1=VALUE_1
KEY_2=VALUE_2

# invalid name
@!!??🙈🙉🙊=abc
`

	// Highest priority: .<environment>.env file.
	testEnvFile := `

# blank files above and below to test skipping

# override KEY_2
KEY_2=VALUE_2_OVERRIDE    # Ignore inline comments
`
	// Create temp working directory.
	tempDir := t.TempDir()
	os.Setenv("TAGDB_ENV", "test")

	oldCwd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldCwd)

	// Write env files.
	if err := os.WriteFile(".env", []byte(envFile), 0644); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	if err := os.WriteFile(".test.env", []byte(testEnvFile), 0644); err != nil {
		t.Fatalf("failed to write .test.env file: %v", err)
	}

	// Act.
	loadEnvFiles(tempDir)

	// Assert.
	expectedValues := map[string]string{
		"KEY_0": "VALUE_0",          // from OS env
		"KEY_1": "VALUE_1",          // from .env
		"KEY_2": "VALUE_2_OVERRIDE", // from .test.env
	}

	for k, v := range expectedValues {
		actual := os.Getenv(k)
		if actual != v {
			t.Errorf("unexpected value for env var `%s`: expected `%s`, got `%s`", k, v, actual)
		}
	}
}
