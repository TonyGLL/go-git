package gogit

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
		line := scanner.Text()

		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			log.Printf("Skipping line with incorrect format: %s", line)
			continue
		}

		header := strings.Fields(parts[0])
		if len(header) != 3 {
			log.Printf("Skipping line with incorrect format: %s", line)
			continue
		}

		hash := header[2]
		path := parts[1]
		treeMap[path] = hash
	}

	// 8. Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return treeMap, fmt.Errorf("error scanning tree file: %w", err)
	}

	return treeMap, nil
}
