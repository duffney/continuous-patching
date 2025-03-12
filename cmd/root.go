package cmd

/*
copamatic registry <url>|<userName> <token> --list
copamatic registry <url>|<userName> <token> --list -o matrix.json
copamatic registry patch
*/

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "copamatic",
	Short: "A CLI tool for patching container images",
	Long: `Copamatic is a command line tool that helps you patch container images.
It allows you to easily apply patches to your container images and manage them effectively.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(registryCmd)
}
