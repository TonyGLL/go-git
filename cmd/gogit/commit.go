package gogit

import (
	"fmt"
	"os"

	"github.com/TonyGLL/go-git/internal/gogit"
	"github.com/spf13/cobra"
)

var commitMessage string
var addCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Add commit message to gogit repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := gogit.AddCommit(&commitMessage); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCommitCmd)
	addCommitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message (mandatory)")
	addCommitCmd.MarkFlagRequired("message")
}
