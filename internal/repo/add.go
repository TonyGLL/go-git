package repo

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	repoPath    = filepath.Join(".", ".gogit")
	objectsPath = filepath.Join(repoPath, "objects")
	indexPath   = filepath.Join(repoPath, "index")
)

func AddFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	blobHash, buffer, err := hasObject(content)
	if err != nil {
		return fmt.Errorf("error hashing file: %w", err)
	}

	firstTwo := blobHash[:2]
	rest := blobHash[2:]

	// Create necessary directories
	dirs := []string{firstTwo}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(objectsPath, dir), 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	// Create initial files like newFile
	newFilePath := filepath.Join(fmt.Sprintf("%s/%s", objectsPath, firstTwo), rest)
	if err := os.WriteFile(newFilePath, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("error creating HEAD file: %w", err)
	}

	if err := searchRewriteIndex(blobHash, filePath); err != nil {
		return fmt.Errorf("error updating INDEX file: %w", err)
	}

	return nil
}

func hasObject(content []byte) (string, bytes.Buffer, error) {
	// 1. Create a buffer to build the Git blob object.
	var buffer bytes.Buffer
	// 2. Write the blob header, including the type ("blob"), a space and the length of the content.
	// The `len()` function in Go returns the number of bytes in the slice.
	buffer.WriteString(fmt.Sprintf("blob %d", len(content)))
	// 3. Write the null byte ('\0'), which separates the header from the content.
	buffer.WriteByte(0)
	// 4. Write the actual content (the `[]byte`).
	buffer.Write(content)
	// 5. Compute the SHA-1 hash of the entire byte sequence.
	hash := sha1.Sum(buffer.Bytes())
	// 6. Format the resulting hash as a hexadecimal string.
	blobHash := fmt.Sprintf("%x", hash)
	return blobHash, buffer, nil
}

func searchRewriteIndex(blobHash string, filePath string) error {
	// 1. Open the index file for reading
	indexContent, err := os.Open(indexPath)
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
		return fmt.Errorf("error reading file %s: %w", indexPath, err)
	}

	// 4. If after reading the file we didn't find the filepath, append it.
	if !found {
		newLine := fmt.Sprintf("%s %s", blobHash, filePath)
		lines = append(lines, newLine)
	}

	// 5. Join all lines into a single string and overwrite the file.
	output := strings.Join(lines, "\n")
	// Add a final newline so the file ends correctly.
	err = os.WriteFile(indexPath, []byte(output+"\n"), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", indexPath, err)
	}

	fmt.Printf("File '%s' updated successfully.\n", indexPath)
	return nil
}
