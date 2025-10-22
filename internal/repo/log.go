package repo

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TonyGLL/go-git/pkg"
)

func LogRepo() error {
	headFile, err := os.Open(pkg.HeadPath)
	if err != nil {
		return err
	}
	defer headFile.Close()

	var headRef string
	headScanner := bufio.NewScanner(headFile)
	for headScanner.Scan() {
		line := headScanner.Text()

		words := strings.Fields(line)
		headRef = words[1]
	}

	branchRefPath := fmt.Sprintf("%s/%s", pkg.RepoPath, headRef)
	branchHashFile, err := os.Open(branchRefPath)
	if err != nil {
		return err
	}
	defer branchHashFile.Close()

	var currentHash string
	brandHashScanner := bufio.NewScanner(branchHashFile)
	for brandHashScanner.Scan() {
		currentHash = brandHashScanner.Text()
	}

	err = pkg.ReadObject(currentHash)
	if err != nil {
		return err
	}

	return nil
}
