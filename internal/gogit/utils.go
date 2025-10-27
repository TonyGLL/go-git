package gogit

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

// BuildWorkdirMap walks the repoRoot and returns a map of relative path -> sha1hex.
func BuildWorkdirMap() (map[string]string, error) {
	repoRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get the current directory: %v", err)
	}
	workdirMap := make(map[string]string)

	// 1. Load the .gogitignore rules.
	ignorePatterns, err := parseGitignore(repoRoot)
	if err != nil {
		return nil, err
	}

	// 2. Start the recursive walk.
	walkErr := filepath.WalkDir(repoRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path to the repository root.
		relativePath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		// Normalize to forward slashes for consistent comparison.
		relativePath = filepath.ToSlash(relativePath)

		// Skip the root (".").
		if relativePath == "." {
			return nil
		}

		// Ignore the .gogitignore file itself.
		if !d.IsDir() && filepath.Base(relativePath) == ".gogitignore" {
			return nil
		}

		// Strict Filter: always ignore the .gogit directory.
		if d.IsDir() && (relativePath == ".gogit" || strings.HasPrefix(relativePath, ".gogit/")) {
			return filepath.SkipDir
		}

		// Evaluate ignore rules (they are applied in order; '!' negations undo previous ignores).
		isIgnored := false
		name := d.Name() // name of the file/dir
		isDir := d.IsDir()

		for _, rawPattern := range ignorePatterns {
			if rawPattern == "" {
				continue
			}
			pattern := filepath.ToSlash(strings.TrimSpace(rawPattern))

			negated := false
			if strings.HasPrefix(pattern, "!") {
				negated = true
				pattern = strings.TrimPrefix(pattern, "!")
				pattern = strings.TrimSpace(pattern)
				if pattern == "" {
					// invalid "!" pattern -> ignore
					continue
				}
			}

			// If the pattern ends in "/" it points to directories.
			patternDirOnly := strings.HasSuffix(pattern, "/")
			if patternDirOnly {
				pattern = strings.TrimSuffix(pattern, "/")
			}

			matched := false

			// If the pattern contains a '/' we compare it against the full relative path.
			if strings.Contains(pattern, "/") {
				// if pattern starts with "/" we treat it as relative to the root: remove prefix if it exists
				if strings.HasPrefix(pattern, "/") {
					pattern = strings.TrimPrefix(pattern, "/")
				}
				// Match using filepath.Match against relativePath
				if ok, matchErr := filepath.Match(pattern, relativePath); matchErr == nil && ok {
					matched = true
				} else if matchErr != nil {
					// invalid pattern — we ignore it
					continue
				}
			} else {
				// does not contain '/', compare against the name of the file/dir
				if ok, matchErr := filepath.Match(pattern, name); matchErr == nil && ok {
					matched = true
				} else if matchErr != nil {
					continue
				}
			}

			// If the pattern is exclusive to directories, and this is not a dir -> no match
			if matched && patternDirOnly && !isDir {
				matched = false
			}

			if matched {
				if negated {
					// A negation undoes the ignore state.
					isIgnored = false
				} else {
					isIgnored = true
				}
				// We don't break; git processes all lines (last relevant match).
			}
		}

		// If it's ignored -> if it's a dir, avoid entering; if it's a file, skip it.
		if isIgnored {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// If it's a valid directory, we do nothing (only files are hashed).
		if d.IsDir() {
			return nil
		}

		// 3. Process and hash each valid file.
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read the file %s: %w", path, err)
		}

		// Create a "blob <size>\0" header.
		header := fmt.Sprintf("blob %d\x00", len(content))

		// Concatenate and hash.
		hasher := sha1.New()
		_, _ = hasher.Write([]byte(header))
		_, _ = hasher.Write(content)
		hashBytes := hasher.Sum(nil)

		// Convert to hexadecimal and store.
		hashHex := hex.EncodeToString(hashBytes)
		// Save with the relative path (without "./").
		workdirMap[relativePath] = hashHex

		return nil
	})

	if walkErr != nil {
		return nil, fmt.Errorf("error during the directory walk: %w", walkErr)
	}

	return workdirMap, nil
}

// parseGitignore reads .gogitignore and returns the lines in order (including negations).
func parseGitignore(repoRoot string) ([]string, error) {
	ignoreFilePath := filepath.Join(repoRoot, ".gogitignore")

	content, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		// If there is no .gogitignore, it is not an error, there is simply nothing to ignore.
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("error reading .gogitignore: %w", err)
	}

	var patterns []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Ignore empty lines and comments.
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}
		patterns = append(patterns, trimmedLine)
	}
	return patterns, nil
}
