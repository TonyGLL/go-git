package pkg

import (
	"bufio"
	"fmt"
	"os"
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
