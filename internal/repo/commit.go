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

func AddCommit(message *string) error {
	hashPathMap := make(map[string]string)

	file, err := os.Open(pkg.IndexPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 3. Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// 4. Iterate over each line of the file
	for scanner.Scan() {
		line := scanner.Text() // Get the line as a string

		// 5. Split the line into a slice of words
		words := strings.Fields(line)

		// 6. Check that there are at least two words
		if len(words) < 2 {
			log.Printf("Skipping line with incorrect format 0: %s", line)
			continue // Go to the next line if the format is incorrect
		}

		// 7. Assign the first word as key and the second as value
		key := words[0]
		value := words[1]
		hashPathMap[key] = value
	}

	if len(hashPathMap) < 1 {
		log.Printf("No file to commit")
		return nil
	}

	// 8. Check for errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	headRef := make(map[string]string)
	ref, err := os.Open(pkg.HeadPath)
	if err != nil {
		return err
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
			log.Printf("Skipping line with incorrect format 1: %s", line)
			continue // Go to the next line if the format is incorrect
		}

		key := words[0]
		value := words[1]
		headRef[key] = value
	}

	// 8. Check for errors during scanning
	if err := scannerRef.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	var parentCommitHash string
	branchRef, err := os.Open(fmt.Sprintf("%s/%s", pkg.RepoPath, headRef["ref:"]))
	if err != nil {
		return err
	}
	defer branchRef.Close()

	// 3. Create a scanner to read the file line by line
	scannerBranchRef := bufio.NewScanner(branchRef)

	// 4. Iterate over each line of the file
	for scannerBranchRef.Scan() {
		parentCommitHash = scannerBranchRef.Text()
	}

	// 8. Check for errors during scanning
	if err := scannerBranchRef.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	authorName := "TonyGLL"
	commitHash, commitContent, err := pkg.HashCommit(parentCommitHash, authorName, *message, hashPathMap)
	if err != nil {
		return fmt.Errorf("error hashing file: %w", err)
	}

	firstTwo := commitHash[:2]
	rest := commitHash[2:]

	// Create necessary directories
	dirs := []string{firstTwo}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(pkg.ObjectsPath, dir), 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	// Create initial files like newFile
	newFilePath := filepath.Join(fmt.Sprintf("%s/%s", pkg.ObjectsPath, firstTwo), rest)
	if err := os.WriteFile(newFilePath, commitContent, 0644); err != nil {
		return fmt.Errorf("error creating Object file: %w", err)
	}

	// Create initial files like newFile
	newRefHeadPath := filepath.Join(fmt.Sprintf("%s/%s", pkg.RefHeadsPath, "main"))
	if err := os.WriteFile(newRefHeadPath, []byte(commitHash+"\n"), 0644); err != nil {
		return fmt.Errorf("error creating Object file: %w", err)
	}

	err = os.WriteFile(pkg.IndexPath, []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", pkg.IndexPath, err)
	}

	return nil
}
