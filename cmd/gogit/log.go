package gogit

import (
	"fmt"
	"os"

	"github.com/TonyGLL/go-git/internal/gogit"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commits logs",
	Run: func(cmd *cobra.Command, args []string) {
		if err := gogit.LogRepo(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(logCmd)
}
