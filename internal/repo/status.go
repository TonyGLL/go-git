package repo

import (
	"fmt"

	"github.com/TonyGLL/go-git/pkg"
)

func StatusRepo() error {
	currentHash, err := pkg.GetBranchHash()
	if err != nil {
		return err
	}

	var treeMap map[string]string
	if currentHash != "" {
		lastCommit, err := pkg.ReadCommit(currentHash)
		if err != nil {
			return err
		}

		lastTreeHash := lastCommit.Tree
		treeMap, err = pkg.ReadTree(lastTreeHash)
		if err != nil {
			return err
		}
	}

	indexMap, err := pkg.ReadIndex()
	if err != nil {
		return err
	}
	statusInfo := &pkg.StatusInfo{
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

	workdirMap, err := pkg.BuildWorkdirMap()
	if err != nil {
		return fmt.Errorf("no se pudo construir el mapa del directorio de trabajo: %w", err)
	}

	for path, workdirHash := range workdirMap {
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

	pkg.PrintStatus(statusInfo)

	return nil
}
