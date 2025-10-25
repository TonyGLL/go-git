package gogit

func LogRepo() error {
	currentHash, err := GetBranchHash()
	if err != nil {
		return err
	}

	err = ReadObject(currentHash)
	if err != nil {
		return err
	}

	return nil
}
