package commands

import (
	"fmt"
	"os"

	"github.com/TonyGLL/go-git/internal/repo"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <file|directory>",
	Short: "Add a file or directory to the gogit repository",
	Long: `Adds the specified file or directory to the staging area (index).
When a directory is specified, it recursively adds all files within that
directory, excluding the .gogit directory itself.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pathToAdd := args[0]

		if err := repo.Add(pathToAdd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("File(s) added successfully.")
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
