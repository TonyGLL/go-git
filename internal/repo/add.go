package repo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TonyGLL/go-git/pkg"
)

func AddFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	blobHash, buffer, err := pkg.HashObject(content)
	if err != nil {
		return fmt.Errorf("error hashing file: %w", err)
	}

	firstTwo := blobHash[:2]
	rest := blobHash[2:]
	objectPath := filepath.Join(pkg.ObjectsPath, firstTwo, rest)

	// Check if object exists. If not, create it.
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		// Create necessary directories
		if err := os.MkdirAll(filepath.Join(pkg.ObjectsPath, firstTwo), 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", filepath.Join(pkg.ObjectsPath, firstTwo), err)
		}

		// Write the blob object content
		if err := os.WriteFile(objectPath, buffer.Bytes(), 0644); err != nil {
			return fmt.Errorf("error writing blob object to %s: %w", objectPath, err)
		}
	} else if err != nil {
		// An error other than "not exist" occurred
		return fmt.Errorf("error checking object existence at %s: %w", objectPath, err)
	}

	// Always update the index with the file and its hash
	if err := searchRewriteIndex(blobHash, filePath); err != nil {
		return fmt.Errorf("error updating INDEX file: %w", err)
	}

	return nil
}

func searchRewriteIndex(blobHash string, filePath string) error {
	// 1. Open the index file for reading
	indexContent, err := os.Open(pkg.IndexPath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var lines []string     // A slice to store the new contents of the file.
	var found bool = false // Flag to indicate if the filepath was found.

	scanner := bufio.NewScanner(indexContent)
	// 2. Read the file line by line.
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into two parts at the first space.
		// Use SplitN to handle file paths that may contain spaces.
		parts := strings.SplitN(line, " ", 2)

		// Verify that the line has the expected format (blobhash filepath).
		if len(parts) == 2 {
			currentFilepath := parts[1]

			// 3. Check if the filepath matches our target.
			if currentFilepath == filePath {
				// Match! Update the line with the new blob hash.
				newLine := fmt.Sprintf("%s %s", blobHash, filePath)
				lines = append(lines, newLine)
				found = true // Mark that we found it.
			} else {
				// No match, keep the original line.
				lines = append(lines, line)
			}
		} else {
			// If the line doesn't have the expected format (or is empty), keep it.
			lines = append(lines, line)
		}
	}

	// Check for errors that occurred during scanning.
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", pkg.IndexPath, err)
	}

	// 4. If after reading the file we didn't find the filepath, append it.
	if !found {
		newLine := fmt.Sprintf("%s %s", blobHash, filePath)
		lines = append(lines, newLine)
	}

	// 5. Join all lines into a single string and overwrite the file.
	output := strings.Join(lines, "\n")
	// Add a final newline so the file ends correctly.
	err = os.WriteFile(pkg.IndexPath, []byte(output+"\n"), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", pkg.IndexPath, err)
	}

	fmt.Printf("File '%s' updated successfully.\n", pkg.IndexPath)
	return nil
}
