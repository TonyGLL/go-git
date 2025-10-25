package gogit

import (
	"fmt"
	"os"

	"github.com/TonyGLL/go-git/internal/gogit"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show commit status",
	Run: func(cmd *cobra.Command, args []string) {
		if err := gogit.StatusRepo(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
