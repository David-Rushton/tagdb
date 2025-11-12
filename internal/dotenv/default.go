/*
Reads .env files from the current working directory.

Precedence:
1. OS environment variables
2. .env file
3. .<environment>.env file

Where <environment> is defined via a TAGDB_ENV environment variable.

Env takes a best-effort approach.  Failures are logged, but do not prevent
application startup.
*/
package dotenv

import (
	"bufio"
	"io"
	"os"
	"slices"
	"strings"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

const (
	tagDbEnv         = "TAGDB_ENV"
	envFileExtension = ".env"
)

func init() {
	// CWD.
	loadEnvFiles(".")
}

func loadEnvFiles(root string) {
	// Search for env files in current working directory.
	entries, err := os.ReadDir(root)
	if err != nil {
		logger.Warnf("cannot read env files from current working directory because %s", err)
		return
	}

	// Get or create env.
	environment := os.Getenv(tagDbEnv)
	if environment == "" {
		environment = "development"
		os.Setenv(tagDbEnv, environment)
	}

	envFilesToLoad := []string{envFileExtension, "." + environment + envFileExtension}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// os.ReadDir always returns files in alphabetical order.
		// Meaning:
		//  - OS env values are overridden by .env file values.
		//  - .env file values are overridden by .<environment>.env file values.
		if slices.Contains(envFilesToLoad, entry.Name()) {
			values := parseEnvFile(entry.Name())
			for k, v := range values {
				if err := os.Setenv(k, v); err != nil {
					logger.Warnf("cannot set env var `%s` from file `%s` to value  `%s` because %s", k, entry.Name(), v, err)
				}
			}
		}
	}
}

func parseEnvFile(path string) map[string]string {
	file, err := os.Open(path)
	if err != nil {
		logger.Warnf("cannot open env file `%s` because %s", path, err)
		return nil
	}
	defer file.Close()

	r := bufio.NewReader(file)
	result := make(map[string]string)

	for {
		// Read next line.
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			logger.Warnf("cannot read next line of env files `%s` because `%s`", path, err)
			continue
		}

		line = strings.TrimSpace(line)

		// Skip comments and empty lines.
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split key and value.
		name, rawValue, found := strings.Cut(line, "=")
		value, _, _ := strings.Cut(rawValue, "#")
		if !found {
			logger.Warnf("cannot parse malformed environment variable `%s`", line)
			continue
		}

		result[strings.TrimSpace(name)] = strings.TrimSpace(value)
	}

	return result
}
