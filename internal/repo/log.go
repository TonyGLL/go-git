package repo

import (
	"github.com/TonyGLL/go-git/pkg"
)

func LogRepo() error {
	currentHash, err := pkg.GetBranchHash()
	if err != nil {
		return err
	}

	err = pkg.ReadObject(currentHash)
	if err != nil {
		return err
	}

	return nil
}
