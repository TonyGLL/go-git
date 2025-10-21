package cmd

import (
	"fmt"
	"os"

	"github.com/TonyGLL/go-git/internal/repo"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gogit", // The name of your application
	Short: "gogit - a simplified Git replica written in Go",
	Long: `gogit is a minimalist version control system
created as a learning project to understand the fundamental
concepts of Git.`,
	// The code here will be executed if 'gogit' is called without any subcommands.
	// Typically, it displays help. Cobra does this by default.
}

// initCmd for 'gogit init'
var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Creates a new gogit repository",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}

		// Here's the magic: the CLI calls the internal logic
		if err := repo.InitRepo(targetDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var addFileCmd = &cobra.Command{
	Use:   "add <file>",
	Short: "Add file to gogit repository",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := "."
		if len(args) > 0 {
			filePath = args[0]
		}

		// Here's the magic: the CLI calls the internal logic
		if err := repo.AddFile(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var commitMessage string
var addCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Add commit message to gogit repository",
	Run: func(cmd *cobra.Command, args []string) {
		// Here's the magic: the CLI calls the internal logic
		if err := repo.AddCommit(&commitMessage); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is the main function that calls main.main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initCmd for 'gogit init'
func init() {
	// Add subcommands to the root command
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addFileCmd)
	rootCmd.AddCommand(addCommitCmd)
	addCommitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message (mandatory)")
	addCommitCmd.MarkFlagRequired("message")
}
