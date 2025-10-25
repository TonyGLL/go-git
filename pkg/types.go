package pkg

import "time"

// Commit represents a commit object.
type Commit struct {
	Hash    string
	Tree    string
	Parent  string
	Author  string
	Date    time.Time
	Message string
}

type StatusInfo struct {
	Branch    string
	Staged    []string
	Unstaged  []string
	Untracked []string
}
