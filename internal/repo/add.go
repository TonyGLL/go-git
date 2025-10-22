package repo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TonyGLL/go-git/pkg"
)

// Add handles adding files to the repository's index.
func Add(path string) error {
	// 1. Read the index file once into memory.
	indexEntries, err := readIndex()
	if err != nil {
		return fmt.Errorf("error reading index: %w", err)
	}

	if path == "." {
		// Walk the current directory
		err = filepath.Walk(".", func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Ignore the .gogit directory
			if info.IsDir() && info.Name() == ".gogit" {
				return filepath.SkipDir
			}
			// Ignore other directories
			if info.IsDir() {
				return nil
			}

			// Normalize path for consistent checks
			normalizedPath := filepath.ToSlash(filePath)
			if strings.HasPrefix(normalizedPath, ".gogit/") {
				return nil
			}

			// 2. Process each file and update the in-memory map
			fmt.Printf("Adding '%s'\n", filePath)
			return processFile(filePath, indexEntries)
		})
		if err != nil {
			return err
		}
	} else {
		// If it's not ".", treat it as a single file or directory
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("error stating path %s: %w", path, err)
		}
		if info.IsDir() {
			return fmt.Errorf("adding single directories is not supported, use 'add .' instead")
		}

		if err := processFile(path, indexEntries); err != nil {
			return err
		}
	}

	// 3. Write the updated index back to the file once.
	if err := writeIndex(indexEntries); err != nil {
		return fmt.Errorf("error writing index file: %w", err)
	}

	return nil
}

// processFile handles hashing a single file and adding it to the in-memory index map.
func processFile(filePath string, indexEntries map[string]string) error {
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
		if err := os.MkdirAll(filepath.Join(pkg.ObjectsPath, firstTwo), 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", filepath.Join(pkg.ObjectsPath, firstTwo), err)
		}
		if err := os.WriteFile(objectPath, buffer.Bytes(), 0644); err != nil {
			return fmt.Errorf("error writing blob object to %s: %w", objectPath, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking object existence at %s: %w", objectPath, err)
	}

	// Update the in-memory map
	indexEntries[filePath] = blobHash
	return nil
}

// readIndex reads the index file into a map.
func readIndex() (map[string]string, error) {
	indexEntries := make(map[string]string)
	indexFile, err := os.Open(pkg.IndexPath)
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
func writeIndex(indexEntries map[string]string) error {
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

	if err := os.WriteFile(pkg.IndexPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("error writing to index file %s: %w", pkg.IndexPath, err)
	}
	return nil
}
