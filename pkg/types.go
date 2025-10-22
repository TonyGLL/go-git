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
