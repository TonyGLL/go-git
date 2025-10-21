package pkg

import "path/filepath"

var (
	RepoPath     = filepath.Join(".", ".gogit")
	ObjectsPath  = filepath.Join(RepoPath, "objects")
	IndexPath    = filepath.Join(RepoPath, "index")
	HeadPath     = filepath.Join(RepoPath, "HEAD")
	RefHeadsPath = filepath.Join(RepoPath, "refs/heads")
)
