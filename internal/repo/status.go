package repo

import (
	"github.com/TonyGLL/go-git/pkg"
)

func StatusRepo() error {
	currentHash, err := pkg.GetBranchHash()
	if err != nil {
		return err
	}

	lastCommit, err := pkg.ReadCommit(currentHash)
	if err != nil {
		return err
	}

	lastTreeHash := lastCommit.Tree
	treeMap, err := pkg.ReadTree(lastTreeHash)
	if err != nil {
		return err
	}

	indexMap, err := pkg.ReadIndex()
	if err != nil {
		return err
	}

	pkg.PrintStatus(&pkg.StatusInfo{
		Branch:    "main",
		Staged:    []string{"aaa"},
		Unstaged:  []string{"aaa"},
		Untracked: []string{"aaa"},
	})

	return nil
}
