package gogit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gogit",
	Short: "gogit - a simplified Git replica written in Go",
	Long: `gogit is a minimalist version control system
created as a learning project to understand the fundamental
concepts of Git.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
