package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func GetHeadRef() (map[string]string, error) {
	headRef := make(map[string]string)
	ref, err := os.Open(HeadPath)
	if err != nil {
		return nil, err
	}
	defer ref.Close()

	// 3. Create a scanner to read the file line by line
	scannerRef := bufio.NewScanner(ref)

	// 4. Iterate over each line of the file
	for scannerRef.Scan() {
		line := scannerRef.Text() // Get the line as a string

		// 5. Split the line into a slice of words
		words := strings.Fields(line)

		// 6. Check that there are at least two words
		if len(words) < 2 {
			log.Printf("Skipping line with incorrect format: %s", line)
			continue // Go to the next line if the format is incorrect
		}

		key := words[0]
		value := words[1]
		headRef[key] = value
	}

	// 8. Check for errors during scanning
	if err := scannerRef.Err(); err != nil {
		return nil, fmt.Errorf("error scanning HEAD file: %w", err)
	}

	return headRef, nil
}

// readIndex reads the index file into a map.
func ReadIndex() (map[string]string, error) {
	indexEntries := make(map[string]string)
	indexFile, err := os.Open(IndexPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, return an empty map. It will be created on write.
			return indexEntries, nil
		}
		return nil, fmt.Errorf("error opening index for reading: %w", err)
	}
	defer indexFile.Close()

	scanner := bufio.NewScanner(indexFile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			indexEntries[parts[1]] = parts[0] // map[filepath] = hash
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning index file: %w", err)
	}
	return indexEntries, nil
}

// writeIndex writes the map of entries to the index file.
func WriteIndex(indexEntries map[string]string) error {
	var lines []string
	// For deterministic output, sort the file paths before writing.
	var paths []string
	for path := range indexEntries {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		lines = append(lines, fmt.Sprintf("%s %s", indexEntries[path], path))
	}

	output := strings.Join(lines, "\n")
	if len(lines) > 0 {
		output += "\n" // Add a final newline
	}

	if err := os.WriteFile(IndexPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("error writing to index file %s: %w", IndexPath, err)
	}
	return nil
}

func GetBranchHash() (string, error) {
	headFile, err := os.Open(HeadPath)
	if err != nil {
		return "", err
	}
	defer headFile.Close()

	var headRef string
	headScanner := bufio.NewScanner(headFile)
	for headScanner.Scan() {
		line := headScanner.Text()

		words := strings.Fields(line)
		headRef = words[1]
	}

	branchRefPath := fmt.Sprintf("%s/%s", RepoPath, headRef)
	branchHashFile, err := os.Open(branchRefPath)
	if err != nil {
		return "", err
	}
	defer branchHashFile.Close()

	var currentHash string
	brandHashScanner := bufio.NewScanner(branchHashFile)
	for brandHashScanner.Scan() {
		currentHash = brandHashScanner.Text()
	}

	return currentHash, nil
}
