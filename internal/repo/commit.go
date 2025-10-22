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
		return fmt.Errorf("error scanning index file: %w", err)
	}

	// --- Generate and save the Tree object ---
	treeHash, treeContent, err := pkg.HashTree(hashPathMap)
	if err != nil {
		return fmt.Errorf("error hashing tree: %w", err)
	}

	treeFirstTwo := treeHash[:2]
	treeRest := treeHash[2:]

	// Create necessary directories for the tree object
	if err := os.MkdirAll(filepath.Join(pkg.ObjectsPath, treeFirstTwo), 0755); err != nil {
		return fmt.Errorf("error creating directory for tree object %s: %w", treeFirstTwo, err)
	}

	// Write the tree object content
	treeObjectPath := filepath.Join(pkg.ObjectsPath, treeFirstTwo, treeRest)
	if err := os.WriteFile(treeObjectPath, treeContent, 0644); err != nil {
		return fmt.Errorf("error writing tree object to %s: %w", treeObjectPath, err)
	}
	// --- End Tree object generation ---


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
		return fmt.Errorf("error scanning HEAD file: %w", err)
	}

	var parentCommitHash string
	branchRefPath := fmt.Sprintf("%s/%s", pkg.RepoPath, headRef["ref:"])
	branchRef, err := os.Open(branchRefPath)
	if err != nil {
		// If the branch ref file doesn't exist, it's likely the first commit.
		// In this case, parentCommitHash remains empty.
		if os.IsNotExist(err) {
			parentCommitHash = ""
		} else {
			return err
		}
	} else {
		defer branchRef.Close()
		scannerBranchRef := bufio.NewScanner(branchRef)
		for scannerBranchRef.Scan() {
			parentCommitHash = scannerBranchRef.Text()
		}
		if err := scannerBranchRef.Err(); err != nil {
			return fmt.Errorf("error scanning branch ref file: %w", err)
		}
	}


	authorName := "TonyGLL"
	// Call HashCommit with the treeHash
	commitHash, commitContent, err := pkg.HashCommit(treeHash, parentCommitHash, authorName, *message)
	if err != nil {
		return fmt.Errorf("error hashing commit: %w", err)
	}

	firstTwo := commitHash[:2]
	rest := commitHash[2:]

	// Create necessary directories for the commit object
	if err := os.MkdirAll(filepath.Join(pkg.ObjectsPath, firstTwo), 0755); err != nil {
		return fmt.Errorf("error creating directory for commit object %s: %w", firstTwo, err)
	}

	// Create commit object file
	newCommitObjectPath := filepath.Join(pkg.ObjectsPath, firstTwo, rest)
	if err := os.WriteFile(newCommitObjectPath, commitContent, 0644); err != nil {
		return fmt.Errorf("error creating commit object file: %w", err)
	}

	// Update branch reference (e.g., refs/heads/main)
	newRefHeadPath := filepath.Join(pkg.RepoPath, headRef["ref:"])
	if err := os.WriteFile(newRefHeadPath, []byte(commitHash+"\n"), 0644); err != nil {
		return fmt.Errorf("error updating branch reference file: %w", err)
	}

	// Clear the index after a successful commit
	err = os.WriteFile(pkg.IndexPath, []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("error clearing index file: %w", err)
	}

	return nil
}
