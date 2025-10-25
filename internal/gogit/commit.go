package gogit

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ReadCommit reads a commit object from the repository and returns a Commit struct.
func ReadCommit(hash string) (*Commit, error) {
	firstTwo := hash[:2]
	rest := hash[2:]

	currentObjectPath := fmt.Sprintf("%s/%s/%s", ObjectsPath, firstTwo, rest)

	file, err := os.Open(currentObjectPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commit Commit
	commit.Hash = hash

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "tree ") {
			commit.Tree = strings.TrimSpace(strings.TrimPrefix(line, "tree "))
		} else if strings.HasPrefix(line, "parent ") {
			commit.Parent = strings.TrimSpace(strings.TrimPrefix(line, "parent "))
		} else if strings.HasPrefix(line, "author ") {
			commit.Author = strings.TrimSpace(strings.TrimPrefix(line, "author "))
		} else if strings.HasPrefix(line, "date ") {
			dateStr := strings.TrimSpace(strings.TrimPrefix(line, "date "))
			commit.Date, _ = time.Parse(time.RFC3339, dateStr)
		} else if line == "" {
			break // End of headers
		}
	}

	for scanner.Scan() {
		commit.Message += scanner.Text() + "\n"
	}
	commit.Message = strings.TrimSpace(commit.Message)

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &commit, nil
}

func AddCommit(message *string) error {
	indexMap, err := ReadIndex()
	if err != nil {
		return err
	}

	if len(indexMap) < 1 {
		log.Printf("No file to commit")
		return nil
	}

	// --- Generate and save the Tree object ---
	treeHash, treeContent, err := HashTree(indexMap)
	if err != nil {
		return fmt.Errorf("error hashing tree: %w", err)
	}

	treeFirstTwo := treeHash[:2]
	treeRest := treeHash[2:]

	// Create necessary directories for the tree object
	if err := os.MkdirAll(filepath.Join(ObjectsPath, treeFirstTwo), 0755); err != nil {
		return fmt.Errorf("error creating directory for tree object %s: %w", treeFirstTwo, err)
	}

	// Write the tree object content
	treeObjectPath := filepath.Join(ObjectsPath, treeFirstTwo, treeRest)
	if err := os.WriteFile(treeObjectPath, treeContent, 0644); err != nil {
		return fmt.Errorf("error writing tree object to %s: %w", treeObjectPath, err)
	}
	// --- End Tree object generation ---

	headRef, err := GetHeadRef()
	if err != nil {
		return err
	}

	var parentCommitHash string
	branchRefPath := fmt.Sprintf("%s/%s", RepoPath, headRef["ref:"])
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
	commitHash, commitContent, err := HashCommit(treeHash, parentCommitHash, authorName, *message)
	if err != nil {
		return fmt.Errorf("error hashing commit: %w", err)
	}

	firstTwo := commitHash[:2]
	rest := commitHash[2:]

	// Create necessary directories for the commit object
	if err := os.MkdirAll(filepath.Join(ObjectsPath, firstTwo), 0755); err != nil {
		return fmt.Errorf("error creating directory for commit object %s: %w", firstTwo, err)
	}

	// Create commit object file
	newCommitObjectPath := filepath.Join(ObjectsPath, firstTwo, rest)
	if err := os.WriteFile(newCommitObjectPath, commitContent, 0644); err != nil {
		return fmt.Errorf("error creating commit object file: %w", err)
	}

	// Update branch reference (e.g., refs/heads/main)
	newRefHeadPath := filepath.Join(RepoPath, headRef["ref:"])
	if err := os.WriteFile(newRefHeadPath, []byte(commitHash+"\n"), 0644); err != nil {
		return fmt.Errorf("error updating branch reference file: %w", err)
	}

	// Clear the index after a successful commit
	// err = os.WriteFile(IndexPath, []byte(""), 0644)
	// if err != nil {
	// 	return fmt.Errorf("error clearing index file: %w", err)
	// }

	return nil
}
