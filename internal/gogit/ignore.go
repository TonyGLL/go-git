package gogit

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// readGogitignore reads the .gogitignore file and returns a list of patterns.
func readGogitignore() ([]string, error) {
	ignorePath := filepath.Join(".", ".gogitignore")
	file, err := os.Open(ignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // .gogitignore not found, return no patterns
		}
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}

// isIgnored checks if a given path should be ignored based on the patterns
// from .gogitignore.
func isIgnored(path string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}
