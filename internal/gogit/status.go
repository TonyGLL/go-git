package gogit

import (
	"fmt"
)

func StatusRepo() error {
	ignorePatterns, err := readGogitignore()
	if err != nil {
		return fmt.Errorf("error reading .gogitignore: %w", err)
	}

	currentHash, err := GetBranchHash()
	if err != nil {
		return err
	}

	var treeMap map[string]string
	if currentHash != "" {
		lastCommit, err := ReadCommit(currentHash)
		if err != nil {
			return err
		}

		lastTreeHash := lastCommit.Tree
		treeMap, err = ReadTree(lastTreeHash)
		if err != nil {
			return err
		}
	}

	indexMap, err := ReadIndex()
	if err != nil {
		return err
	}
	statusInfo := &StatusInfo{
		Branch:    "main",
		Staged:    []string{},
		Unstaged:  []string{},
		Untracked: []string{},
	}

	for path, indexHash := range indexMap {
		commitHash, existsInCommit := treeMap[path]
		if !existsInCommit {
			statusInfo.Staged = append(statusInfo.Staged, fmt.Sprintf("new file:   %s", path))
		} else if indexHash != commitHash {
			statusInfo.Staged = append(statusInfo.Staged, fmt.Sprintf("modified:   %s", path))
		}
	}
	for path := range treeMap {
		if _, existsInIndex := indexMap[path]; !existsInIndex {
			statusInfo.Staged = append(statusInfo.Staged, fmt.Sprintf("deleted:    %s", path))
		}
	}

	workdirMap, err := BuildWorkdirMap()
	if err != nil {
		return fmt.Errorf("no se pudo construir el mapa del directorio de trabajo: %w", err)
	}

	filteredWorkdirMap := make(map[string]string)
	for path, hash := range workdirMap {
		ignored, err := isIgnored(path, ignorePatterns)
		if err != nil {
			return fmt.Errorf("error checking ignore patterns for %s: %w", path, err)
		}
		if !ignored {
			filteredWorkdirMap[path] = hash
		}
	}

	for path, workdirHash := range filteredWorkdirMap {
		indexHash, existsInIndex := indexMap[path]
		if !existsInIndex {
			// Caso C: Untracked
			statusInfo.Untracked = append(statusInfo.Untracked, path)
		} else if workdirHash != indexHash {
			// Caso D: Modificado Unstaged
			statusInfo.Unstaged = append(statusInfo.Unstaged, fmt.Sprintf("modified:   %s", path))
		}
	}
	// Itera sobre el index para encontrar borrados unstaged
	for path := range indexMap {
		if _, existsInWorkdir := workdirMap[path]; !existsInWorkdir {
			statusInfo.Unstaged = append(statusInfo.Unstaged, fmt.Sprintf("deleted:    %s", path))
		}
	}

	PrintStatus(statusInfo)

	return nil
}
