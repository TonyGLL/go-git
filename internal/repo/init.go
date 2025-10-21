package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitRepo contains the logic to initialize the repository directory structure.
// It receives the path where the repository will be created.
func InitRepo(path string) error {
	// Path to the main repository directory (.gogit or the name you choose)
	repoPath := filepath.Join(path, ".gogit")

	// Check if it already exists
	if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
		return fmt.Errorf("gogit repository already exists in %s", path)
	}

	fmt.Printf("Initializing empty gogit repository in %s\n", repoPath)

	// Create necessary directories
	dirs := []string{"objects", "refs/heads"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(repoPath, dir), 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	// Create initial files like HEAD
	headPath := filepath.Join(repoPath, "HEAD")
	// By default, HEAD points to the 'main' branch (or 'master')
	content := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(headPath, content, 0644); err != nil {
		return fmt.Errorf("error creating HEAD file: %w", err)
	}

	fmt.Println("Repository initialized successfully!")
	return nil
}
