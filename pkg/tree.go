package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func ReadTree(hash string) (map[string]string, error) {
	treeMap := make(map[string]string)
	treePath := fmt.Sprintf("%s/%s/%s", ObjectsPath, hash[:2], hash[2:])

	treeFile, err := os.Open(treePath)
	if err != nil {
		return treeMap, err
	}

	defer treeFile.Close()

	// 3. Create a scanner to read the file line by line
	scanner := bufio.NewScanner(treeFile)

	// 4. Iterate over each line of the file
	for scanner.Scan() {
		line := scanner.Text() // Get the line as a string

		// 5. Split the line into a slice of words
		words := strings.Fields(line)

		// 6. Check that there are at least two words
		if len(words) < 4 {
			log.Printf("Skipping line with incorrect format: %s", line)
			continue // Go to the next line if the format is incorrect
		}

		key := words[2]
		value := words[3]
		treeMap[key] = value
	}

	// 8. Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return treeMap, fmt.Errorf("error scanning HEAD file: %w", err)
	}

	return treeMap, nil
}
