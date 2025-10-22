package pkg

import "fmt"

// PrintCommit prints a commit object with a stylized format.
func PrintCommit(commit *Commit) {
	fmt.Printf("%scommit %s%s\n", ColorYellow, commit.Hash, ColorReset)
	fmt.Printf("Tree: %s\n", commit.Tree)
	if commit.Parent != "" {
		fmt.Printf("%sParent: %s%s\n", ColorRed, commit.Parent, ColorReset)
	}
	fmt.Printf("%sAuthor: %s%s\n", ColorGreen, commit.Author, ColorReset)
	fmt.Printf("%sDate: %s%s\n", ColorBlue, commit.Date.Format("Mon Jan 2 15:04:05 2006 -0700"), ColorReset)
	fmt.Printf("\n\t%s\n\n", commit.Message)
}
